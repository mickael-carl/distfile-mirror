package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ProdriveTechnologies/distfile-mirror/pkg/schema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func downloadAndStoreFile(ctx context.Context, uri string, desiredSha256 *string, uploader *s3manager.Uploader) (string, uint64, error) {
	// Create a temporary file for storing the file to be downloaded.
	tmpfile, err := ioutil.TempFile("", "download")
	if err != nil {
		return "", 0, err
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// Download and store the file into the temporary file.
	// TODO(edsch): Place upperbound on the maximum file size.
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", 0, err
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", 0, err
	}
	fileSize, err := io.Copy(tmpfile, resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", 0, err
	}

	// Compute file checksum and validate it against what is expected.
	hasher := sha256.New()
	if _, err := tmpfile.Seek(0, 0); err != nil {
		return "", 0, err
	}
	if _, err := io.Copy(hasher, tmpfile); err != nil {
		return "", 0, err
	}
	checksum := hex.EncodeToString(hasher.Sum(nil))
	if desiredSha256 != nil && *desiredSha256 != checksum {
		return "", 0, fmt.Errorf("Downloaded copy of %s has checksum %s, whereas %s was expected", uri, checksum, *desiredSha256)
	}

	// Upload file into S3 bucket.
	if _, err := tmpfile.Seek(0, 0); err != nil {
		return "", 0, err
	}
	_, err = uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String("files"),
		Key:    aws.String(fmt.Sprintf("%s|%d", checksum, fileSize)),
		Body:   tmpfile,
	})
	return checksum, uint64(fileSize), err
}

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

	// The container image lacks a functioning temporary directory by default.
	os.Mkdir("/tmp", 0777)

	db, err := gorm.Open("postgres", *dbAddress)
	if err != nil {
		log.Fatal(err)
	}

	s3Session := session.New(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(*s3AccessKeyId, *s3SecretAccessKey, ""),
		Endpoint:         s3Endpoint,
		Region:           s3Region,
		DisableSSL:       s3DisableSsl,
		S3ForcePathStyle: aws.Bool(true),
	})
	s3Uploader := s3manager.NewUploader(s3Session)

	var files []schema.File
	if r := db.Where("present = false").Find(&files); r.Error != nil {
		log.Fatal(r.Error)
	}

	ctx := context.Background()
	for _, file := range files {
		log.Printf("Downloading %s", file.Uri)

		// TODO(edsch): Make timeout configurable.
		ctx, cancel := context.WithTimeout(ctx, time.Hour)
		checksum, fileSize, err := downloadAndStoreFile(ctx, file.Uri, file.Sha256, s3Uploader)
		cancel()
		if err != nil {
			log.Printf("Failed to download and store: %s", err)
			continue
		}

		// Update database entry to prevent successive download.
		if r := db.Model(&schema.File{}).Where("id = ?", file.Id).Updates(schema.File{
			Sha256:  &checksum,
			Size:    &fileSize,
			Present: true,
		}); r.Error != nil {
			log.Printf("Failed to update file entry in database: %s", r.Error)
			continue
		}
	}
}
