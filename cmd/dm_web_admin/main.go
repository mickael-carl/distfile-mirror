package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/ProdriveTechnologies/distfile-mirror/pkg/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// Regular expression of a UUID identifier. Used for URL matching.
	uuidRegex = "[a-f0-9]{8}-[a-f0-9]{4}-[1-5][a-f0-9]{3}-[a-f0-9]{4}-[a-f0-9]{12}"
)

func main() {
	var (
		dbAddress          = flag.String("db.address", "", "Database server address.")
		proxyPublicAddress = flag.String("proxy.public-address", "", "Public address at which the proxy can be contacted.")
	)
	flag.Parse()

	db, err := gorm.Open("postgres", *dbAddress)
	if err != nil {
		log.Fatal(err)
	}

	templates, err := template.ParseGlob("templates/*")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	util.RegisterHealthPage(db, router)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	NewFrontpageService(templates, router, *proxyPublicAddress)
	NewContainerManagementService(db, templates, router)
	NewFileManagementService(db, templates, router, *proxyPublicAddress)
	log.Fatal(http.ListenAndServe(":80", router))
}
