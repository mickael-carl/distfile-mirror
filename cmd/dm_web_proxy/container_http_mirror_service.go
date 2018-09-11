package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/ProdriveTechnologies/distfile-mirror/pkg/schema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
)

var (
	containerPingPattern      = regexp.MustCompile("(.*/)v2/$")
	containerManifestsPattern = regexp.MustCompile("(.*/)v2/(.*)/manifests/(.*)")
	containerBlobsPattern     = regexp.MustCompile("(.*/)v2/(.*)/blobs/(.*)")
)

type containerHttpMirrorService struct {
	scheme   string
	database *gorm.DB
	s3       *s3.S3
	fallback http.Handler
}

func NewContainerHttpMirrorService(scheme string, database *gorm.DB, s3 *s3.S3, fallback http.Handler) http.Handler {
	return &containerHttpMirrorService{
		scheme:   scheme,
		database: database,
		s3:       s3,
		fallback: fallback,
	}
}

func (ms *containerHttpMirrorService) getRegistry(host string, path string) (*schema.ContainerRegistry, error) {
	registryUrl := url.URL{
		Scheme: ms.scheme,
		Host:   host,
		Path:   path,
	}
	var registry schema.ContainerRegistry
	if r := ms.database.Where("uri = ?", registryUrl.String()).Take(&registry); r.Error != nil {
		if r.RecordNotFound() {
			return nil, nil
		}
		return nil, r.Error
	}
	return &registry, nil
}

func (ms *containerHttpMirrorService) getRepository(host string, path string, repositoryName string) (*schema.ContainerRepository, error) {
	registry, err := ms.getRegistry(host, path)
	if registry == nil {
		return nil, err
	}
	var repository schema.ContainerRepository
	if r := ms.database.Where("registry_id = ? AND repository_name = ?", registry.Id, repositoryName).Take(&repository); r.Error != nil {
		if r.RecordNotFound() {
			return nil, nil
		}
		return nil, r.Error
	}
	return &repository, nil
}

func (ms *containerHttpMirrorService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if matches := containerPingPattern.FindStringSubmatch(req.URL.Path); matches != nil {
		// Serve a HTTP 200 response for ping requests for existing repositories.
		registry, err := ms.getRegistry(req.Host, matches[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if registry == nil {
			goto NoMatch
		}

		w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
		return
	} else if matches := containerManifestsPattern.FindStringSubmatch(req.URL.Path); matches != nil {
		// Serve manifest from database.
		repository, err := ms.getRepository(req.Host, matches[1], matches[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if repository == nil {
			goto NoMatch
		}

		var image schema.ContainerImage
		if r := ms.database.Where("repository_id = ? AND digest = ? AND manifest IS NOT NULL", repository.Id, matches[3]).Take(&image); r.Error != nil {
			if r.RecordNotFound() {
				goto NoMatch
			}
			http.Error(w, r.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(*image.Manifest)), 10))
		w.Header().Set("Content-Type", *image.ManifestMediatype)
		w.Write(*image.Manifest)
		return
	} else if matches := containerBlobsPattern.FindStringSubmatch(req.URL.Path); matches != nil {
		// Serve blob from S3. Don't care whether the blob
		// actually corresponds to one of the blobs in the
		// registry. We have a single CAS from which we serve.
		repository, err := ms.getRepository(req.Host, matches[1], matches[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if repository == nil {
			goto NoMatch
		}

		blob, err := ms.s3.GetObjectWithContext(req.Context(), &s3.GetObjectInput{
			Bucket: aws.String("container-blobs"),
			Key:    aws.String(matches[3]),
		})
		if err != nil {
			// TODO(edsch): Return 404 where applicable.
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy blob to HTTP response.
		if blob.ContentLength != nil {
			w.Header().Set("Content-Length", strconv.FormatInt(*blob.ContentLength, 10))
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err = io.Copy(w, blob.Body)
		blob.Body.Close()
		if err != nil {
			log.Printf("Failed to write body: %s", err)
			return
		}
		return
	}

NoMatch:
	ms.fallback.ServeHTTP(w, req)
}
