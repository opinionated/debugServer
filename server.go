package main

import (
	"github.com/gorilla/mux"
	"github.com/opinionated/utils/log"
	"net/http"
)

func startServer() {
	log.Info("starting server")

	http.Handle("/", newHandler())
}

func newHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/add", HandleAddArticle).Methods("POST")
	router.HandleFunc("/frontpage", HandleGetFrontpage).Methods("GET")
	return router
}
