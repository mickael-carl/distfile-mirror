package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ProdriveTechnologies/distfile-mirror/pkg/schema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/docker/distribution/manifest/manifestlist"
	_ "github.com/docker/distribution/manifest/schema1"
	_ "github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/client"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/registry"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	oci_digest "github.com/opencontainers/go-digest"
)

type anonymousCredentialStore struct {
}

func (cs *anonymousCredentialStore) Basic(*url.URL) (string, string) {
	return "", ""
}

func (cs *anonymousCredentialStore) RefreshToken(u *url.URL, service string) string {
	return ""
}

func (cs *anonymousCredentialStore) SetRefreshToken(u *url.URL, service string, token string) {
}

func downloadAndStoreContainerImage(ctx context.Context, registryUrl string, repositoryName string, digest string, uploader *s3manager.Uploader) (string, []byte, error) {
	// Send ping to registry to obtain OAuth2 bearer token.
	parsedRegistryUrl, err := url.Parse(registryUrl)
	if err != nil {
		return "", nil, err
	}
	challengeManager, confirmedV2, err := registry.PingV2Registry(parsedRegistryUrl, http.DefaultTransport)
	if err != nil {
		return "", nil, err
	}
	if !confirmedV2 {
		return "", nil, errors.New("Unsupported registry version")
	}

	// Access repository.
	repositoryRef, err := reference.WithName(repositoryName)
	if err != nil {
		return "", nil, err
	}
	creds := &anonymousCredentialStore{}
	repository, err := client.NewRepository(
		repositoryRef,
		registryUrl,
		transport.NewTransport(
			http.DefaultTransport,
			auth.NewAuthorizer(
				challengeManager,
				auth.NewTokenHandler(
					http.DefaultTransport,
					creds,
					repositoryName,
					"pull"),
				auth.NewBasicHandler(creds))))
	if err != nil {
		return "", nil, err
	}

	// Obtain manifest of digest.
	manifestService, err := repository.Manifests(ctx)
	if err != nil {
		return "", nil, err
	}
	manifest, err := manifestService.Get(ctx, oci_digest.Digest(digest))
	if err != nil {
		return "", nil, err
	}

	// Copy blobs into S3. Don't do this for manifest lists. Those
	// are references to other manifests based on operating system
	// and hardware architecture. Users should choose which specific
	// image they want, as storing all of them will use an excessive
	// amount of space.
	if _, ok := manifest.(*manifestlist.DeserializedManifestList); ok {
		log.Printf("Not downloading any blobs, as %s is a manifest list", digest)
	} else {
		blobsService := repository.Blobs(ctx)
		if err != nil {
			return "", nil, err
		}
		for _, descriptor := range manifest.References() {
			// TODO(edsch): Skip downloading if blob already present in the cache.
			log.Printf("Downloading blob %s", descriptor.Digest)
			r, err := blobsService.Open(ctx, descriptor.Digest)
			if err != nil {
				return "", nil, err
			}
			_, err = uploader.UploadWithContext(ctx, &s3manager.UploadInput{
				Bucket: aws.String("container-blobs"),
				Key:    aws.String(string(descriptor.Digest)),
				Body:   r,
			})
			r.Close()
			if err != nil {
				return "", nil, err
			}
		}
	}
	return manifest.Payload()
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

	var containerImages []schema.ContainerImage
	if r := db.Where("manifest IS NULL").Find(&containerImages); r.Error != nil {
		log.Fatal(r.Error)
	}

	ctx := context.Background()
	for _, containerImage := range containerImages {
		// Obtain full information for container image to download.
		var containerRepository schema.ContainerRepository
		if r := db.Where("id = ?", containerImage.RepositoryId).Take(&containerRepository); r.Error != nil {
			log.Printf("Failed to get container repository %s: %s", containerImage.RepositoryId, r.Error)
			continue
		}
		var containerRegistry schema.ContainerRegistry
		if r := db.Where("id = ?", containerRepository.RegistryId).Take(&containerRegistry); r.Error != nil {
			log.Printf("Failed to get container registry %s: %s", containerRepository.RegistryId, r.Error)
			continue
		}
		log.Printf("Downloading %s %s %s", containerRegistry.Uri, containerRepository.RepositoryName, containerImage.Digest)

		// TODO(edsch): Make timeout configurable.
		ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		manifestMediatype, manifest, err := downloadAndStoreContainerImage(ctx, containerRegistry.Uri, containerRepository.RepositoryName, containerImage.Digest, s3Uploader)
		cancel()
		if err != nil {
			log.Printf("Failed to download and store: %s", err)
			continue
		}

		// Update database entry to prevent successive download.
		if r := db.Model(&schema.ContainerImage{}).Where("id = ?", containerImage.Id).Updates(schema.ContainerImage{
			ManifestMediatype: &manifestMediatype,
			Manifest:          &manifest,
		}); r.Error != nil {
			log.Printf("Failed to update container image entry in database: %s", err)
			continue
		}
	}
}
