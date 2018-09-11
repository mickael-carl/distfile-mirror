package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ProdriveTechnologies/distfile-mirror/pkg/schema"
	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	_ "github.com/docker/distribution/manifest/schema1"
	_ "github.com/docker/distribution/manifest/schema2"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type ContainerManagementService struct {
	database  *gorm.DB
	templates *template.Template
}

func NewContainerManagementService(database *gorm.DB, templates *template.Template, router *mux.Router) *ContainerManagementService {
	ms := &ContainerManagementService{
		database:  database,
		templates: templates,
	}
	router.HandleFunc("/containers/", ms.handleRegistriesList)
	router.HandleFunc("/containers/create", ms.handleCreate)
	router.HandleFunc("/containers/registries/{registry_id:"+uuidRegex+"}", ms.handleRegistryInfo)
	router.HandleFunc("/containers/repositories/{repository_id:"+uuidRegex+"}", ms.handleRepositoryInfo)
	router.HandleFunc("/containers/images/{image_id:"+uuidRegex+"}", ms.handleImageInfo)
	return ms
}

func (ms *ContainerManagementService) handleErrorPage(w http.ResponseWriter, req *http.Request, message string, code int) {
	log.Print(message)
	w.WriteHeader(code)
	if err := ms.templates.ExecuteTemplate(w, "error.html", struct {
		Message string
	}{
		Message: message,
	}); err != nil {
		log.Print(err)
	}
}

func (ms *ContainerManagementService) handleCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()

		// Create registry, repository and image if not yet present.
		// TODO(edsch): Store metadata: who creates the image and for what reason.
		var registry schema.ContainerRegistry
		if r := ms.database.FirstOrCreate(&registry, schema.ContainerRegistry{
			Uri: req.Form.Get("registry"),
		}); r.Error != nil {
			ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
			return
		}
		var repository schema.ContainerRepository
		if r := ms.database.FirstOrCreate(&repository, schema.ContainerRepository{
			RegistryId:     registry.Id,
			RepositoryName: req.Form.Get("repository"),
		}); r.Error != nil {
			ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
			return
		}
		var image schema.ContainerImage
		if r := ms.database.FirstOrCreate(&image, schema.ContainerImage{
			RepositoryId: repository.Id,
			Digest:       req.Form.Get("digest"),
		}); r.Error != nil {
			ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/containers/images/"+image.Id, http.StatusSeeOther)
	} else {
		// Present creation form.
		query := req.URL.Query()
		if err := ms.templates.ExecuteTemplate(w, "containers_create.html", struct {
			Registry   string
			Repository string
			Digest     string
		}{
			Registry:   query.Get("registry"),
			Repository: query.Get("repository"),
			Digest:     query.Get("digest"),
		}); err != nil {
			log.Print(err)
		}
	}
}

func (ms *ContainerManagementService) handleRegistriesList(w http.ResponseWriter, req *http.Request) {
	var registries []schema.ContainerRegistry
	if r := ms.database.Order("uri").Find(&registries); r.Error != nil {
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := ms.templates.ExecuteTemplate(w, "containers_registries_list.html", struct {
		Registries []schema.ContainerRegistry
	}{
		Registries: registries,
	}); err != nil {
		log.Print(err)
	}
}

func (ms *ContainerManagementService) handleRegistryInfo(w http.ResponseWriter, req *http.Request) {
	// Obtain registry information.
	var registry schema.ContainerRegistry
	if r := ms.database.Where("id = ?", mux.Vars(req)["registry_id"]).Take(&registry); r.Error != nil {
		// TODO(edsch): Error code.
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Obtain repositories in registry.
	var repositories []schema.ContainerRepository
	if r := ms.database.Where("registry_id = ?", registry.Id).Order("repository_name").Find(&repositories); r.Error != nil {
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := ms.templates.ExecuteTemplate(w, "containers_registry_info.html", struct {
		Registry     *schema.ContainerRegistry
		Repositories []schema.ContainerRepository
	}{
		Registry:     &registry,
		Repositories: repositories,
	}); err != nil {
		log.Print(err)
	}
}

func (ms *ContainerManagementService) handleRepositoryInfo(w http.ResponseWriter, req *http.Request) {
	// Obtain repository information.
	var repository schema.ContainerRepository
	if r := ms.database.Where("id = ?", mux.Vars(req)["repository_id"]).Take(&repository); r.Error != nil {
		// TODO(edsch): Error code.
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}
	var registry schema.ContainerRegistry
	if r := ms.database.Where("id = ?", repository.RegistryId).Take(&registry); r.Error != nil {
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Obtain images in repository.
	// TODO(edsch): Derive image type from ManifestMediatype.
	var images []schema.ContainerImage
	if r := ms.database.Where("repository_id = ?", repository.Id).Order("digest").Find(&images); r.Error != nil {
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := ms.templates.ExecuteTemplate(w, "containers_repository_info.html", struct {
		Registry   *schema.ContainerRegistry
		Repository *schema.ContainerRepository
		Images     []schema.ContainerImage
	}{
		Registry:   &registry,
		Repository: &repository,
		Images:     images,
	}); err != nil {
		log.Print(err)
	}
}

func (ms *ContainerManagementService) handleImageInfo(w http.ResponseWriter, req *http.Request) {
	// Obtain image information.
	var image schema.ContainerImage
	if r := ms.database.Where("id = ?", mux.Vars(req)["image_id"]).Take(&image); r.Error != nil {
		// TODO(edsch): Error code.
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}
	var repository schema.ContainerRepository
	if r := ms.database.Where("id = ?", image.RepositoryId).Take(&repository); r.Error != nil {
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}
	var registry schema.ContainerRegistry
	if r := ms.database.Where("id = ?", repository.RegistryId).Take(&registry); r.Error != nil {
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Extract information from manifest if downloaded. A container
	// image either has layers that allows it to be run directly, or
	// it is a manifest list. When it is a manifest list, it refers
	// to other images based on hardware architecture and OS.
	type manifestInfo struct {
		manifestlist.ManifestDescriptor
		ImageId *string
	}
	var manifests []manifestInfo
	var layers []distribution.Descriptor
	if image.Manifest != nil {
		manifest, _, err := distribution.UnmarshalManifest(*image.ManifestMediatype, *image.Manifest)
		if err != nil {
			ms.handleErrorPage(w, req, err.Error(), http.StatusInternalServerError)
			return
		}
		if manifestList, ok := manifest.(*manifestlist.DeserializedManifestList); ok {
			// Manifest list. Extract which of the manifests
			// contained within are already present in the system.
			var manifestDigests []string
			for _, manifestDescriptor := range manifestList.Manifests {
				manifestDigests = append(manifestDigests, string(manifestDescriptor.Digest))
			}
			var manifestsPresent []schema.ContainerImage
			if r := ms.database.Where("digest IN (?)", manifestDigests).Find(&manifestsPresent); r.Error != nil {
				ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
				return
			}
			imageIds := map[string]string{}
			for _, manifest := range manifestsPresent {
				imageIds[manifest.Digest] = manifest.Id
			}
			for _, manifestDescriptor := range manifestList.Manifests {
				var imageId *string
				if val, ok := imageIds[string(manifestDescriptor.Digest)]; ok {
					imageId = &val
				}
				manifests = append(manifests, manifestInfo{
					ManifestDescriptor: manifestDescriptor,
					ImageId:            imageId,
				})
			}
		} else {
			// Plain container image with layers.
			layers = manifest.References()
		}
	}

	if err := ms.templates.ExecuteTemplate(w, "containers_image_info.html", struct {
		Registry   *schema.ContainerRegistry
		Repository *schema.ContainerRepository
		Image      *schema.ContainerImage
		Manifests  []manifestInfo
		Layers     []distribution.Descriptor
	}{
		Registry:   &registry,
		Repository: &repository,
		Image:      &image,
		Manifests:  manifests,
		Layers:     layers,
	}); err != nil {
		log.Print(err)
	}
}
