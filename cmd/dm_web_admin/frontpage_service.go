package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type FrontpageService struct {
	templates          *template.Template
	proxyPublicAddress string
}

func NewFrontpageService(templates *template.Template, router *mux.Router, proxyPublicAddress string) *FrontpageService {
	fs := &FrontpageService{
		templates:          templates,
		proxyPublicAddress: proxyPublicAddress,
	}
	router.HandleFunc("/", fs.handleFrontpage)
	return fs
}

func (fs *FrontpageService) handleFrontpage(w http.ResponseWriter, req *http.Request) {
	if err := fs.templates.ExecuteTemplate(w, "frontpage.html", struct {
		ProxyPublicAddress string
	}{
		ProxyPublicAddress: fs.proxyPublicAddress,
	}); err != nil {
		log.Print(err)
	}
}
