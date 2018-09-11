package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	var (
		dbAddress = flag.String("db.address", "", "Database server address.")

		s3AccessKeyId     = flag.String("s3.access-key-id", "", "Access key of the S3 bucket holding distfiles")
		s3DisableSsl      = flag.Bool("s3.disable-ssl", false, "Whether SSL should be disabled for the S3 bucket holding distfiles")
		s3Endpoint        = flag.String("s3.endpoint", "", "Endpoint URL of the S3 bucket holding distfiles")
		s3Region          = flag.String("s3.region", "", "Region of the S3 bucket holding distfiles")
		s3SecretAccessKey = flag.String("s3.secret-access-key", "", "Secret access key of the S3 bucket holding distfiles")
	)
	flag.Parse()

	db, err := gorm.Open("postgres", *dbAddress)
	if err != nil {
		panic(err)
	}

	s3Session := session.New(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(*s3AccessKeyId, *s3SecretAccessKey, ""),
		Endpoint:         s3Endpoint,
		Region:           s3Region,
		DisableSSL:       s3DisableSsl,
		S3ForcePathStyle: aws.Bool(true),
	})
	s3 := s3.New(s3Session)

	// Certificate authority used for generating SSL certificates on
	// the fly to 'man in the middle' incoming connections.
	// TODO(edsch): Should we add an adapter to cache certificates?
	caCertificate, err := ioutil.ReadFile("/ca/tls.crt")
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %s", err)
	}
	caPrivateKey, err := ioutil.ReadFile("/ca/tls.key")
	if err != nil {
		log.Fatalf("Failed to load CA private key: %s", err)
	}
	certificateGenerator, err := NewCertificateGenerator(caCertificate, caPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// HTTPS service.
	httpsListener := NewProxyConnectionListener()
	go func() {
		log.Fatal(http.Serve(
			tls.NewListener(httpsListener, &tls.Config{
				GetCertificate: certificateGenerator.GetCertificate,
			}),
			NewFileHttpMirrorService("https", db, s3,
				NewContainerHttpMirrorService("https", db, s3, http.NotFoundHandler()))))
	}()

	// HTTP proxy frontend.
	// TODO(edsch): Demultiplex connections to the appropriate service
	// based on the port number used in the CONNECT request? FTP support?
	log.Fatal(http.ListenAndServe(
		":80",
		NewProxyConnectionHijacker(
			httpsListener,
			NewFileHttpMirrorService("http", db, s3,
				NewContainerHttpMirrorService("http", db, s3, http.NotFoundHandler())))))
}
