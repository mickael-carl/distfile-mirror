package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ProdriveTechnologies/distfile-mirror/pkg/schema"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/gorm"
)

type fileHttpMirrorService struct {
	scheme   string
	database *gorm.DB
	s3       *s3.S3
	fallback http.Handler
}

func NewFileHttpMirrorService(scheme string, database *gorm.DB, s3 *s3.S3, fallback http.Handler) http.Handler {
	return &fileHttpMirrorService{
		scheme:   scheme,
		database: database,
		s3:       s3,
		fallback: fallback,
	}
}

func (ms *fileHttpMirrorService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Determine which CAS entry to serve.
	url := *req.URL
	url.Scheme = ms.scheme
	url.Host = req.Host
	var file schema.File
	if r := ms.database.Where("uri = ? AND present = true", url.String()).Take(&file); r.Error != nil {
		if r.RecordNotFound() {
			// Unknown URL or not yet present in cache.
			ms.fallback.ServeHTTP(w, req)
			return
		}
		http.Error(w, r.Error.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method != http.MethodGet {
		http.Error(w, "Files may only be downloaded using HTTP GET requests", http.StatusMethodNotAllowed)
		return
	}

	// Get blob from storage.
	blob, err := ms.s3.GetObjectWithContext(req.Context(), &s3.GetObjectInput{
		Bucket: aws.String("files"),
		Key:    aws.String(fmt.Sprintf("%s|%d", *file.Sha256, *file.Size)),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy blob to HTTP response.
	w.Header().Set("Content-Length", strconv.FormatUint(*file.Size, 10))
	w.Header().Set("Content-Type", "application/octet-stream")
	copiedSize, err := io.Copy(w, blob.Body)
	blob.Body.Close()
	if err != nil {
		log.Printf("Failed to write body: %s", err)
		return
	}
	if copiedSize != int64(*file.Size) {
		log.Printf("Returned %d, whereas %d was expected", copiedSize, *file.Size)
		return
	}
}
