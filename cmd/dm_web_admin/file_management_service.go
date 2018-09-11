package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ProdriveTechnologies/distfile-mirror/pkg/schema"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type FileManagementService struct {
	database           *gorm.DB
	templates          *template.Template
	proxyPublicAddress string
}

func NewFileManagementService(database *gorm.DB, templates *template.Template, router *mux.Router, proxyPublicAddress string) *FileManagementService {
	ms := &FileManagementService{
		database:           database,
		templates:          templates,
		proxyPublicAddress: proxyPublicAddress,
	}
	router.HandleFunc("/files/", ms.handleFilesList)
	router.HandleFunc("/files/create", ms.handleCreate)
	router.HandleFunc("/files/{file_id:"+uuidRegex+"}", ms.handleFileInfo)
	return ms
}

func (ms *FileManagementService) handleErrorPage(w http.ResponseWriter, req *http.Request, message string, code int) {
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

func (ms *FileManagementService) handleCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()

		// Create file if not yet present.
		// TODO(edsch): Store metadata: who creates the image and for what reason.
		// TODO(edsch): Allow the user to provide a desired SHA-256 sum.
		var file schema.File
		if r := ms.database.FirstOrCreate(&file, schema.File{
			Uri: req.Form.Get("uri"),
		}); r.Error != nil {
			ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, "/files/"+file.Id, http.StatusSeeOther)
	} else {
		// Present creation form.
		if err := ms.templates.ExecuteTemplate(w, "files_create.html", nil); err != nil {
			log.Print(err)
		}
	}
}

func (ms *FileManagementService) handleFilesList(w http.ResponseWriter, req *http.Request) {
	var files []schema.File
	if r := ms.database.Order("uri").Find(&files); r.Error != nil {
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := ms.templates.ExecuteTemplate(w, "files_list.html", struct {
		Files []schema.File
	}{
		Files: files,
	}); err != nil {
		log.Print(err)
	}
}

func (ms *FileManagementService) handleFileInfo(w http.ResponseWriter, req *http.Request) {
	// Obtain file information.
	var file schema.File
	if r := ms.database.Where("id = ?", mux.Vars(req)["file_id"]).Take(&file); r.Error != nil {
		// TODO(edsch): Error code.
		ms.handleErrorPage(w, req, r.Error.Error(), http.StatusInternalServerError)
		return
	}

	if err := ms.templates.ExecuteTemplate(w, "file_info.html", struct {
		File               *schema.File
		ProxyPublicAddress string
	}{
		File:               &file,
		ProxyPublicAddress: ms.proxyPublicAddress,
	}); err != nil {
		log.Print(err)
	}
}
