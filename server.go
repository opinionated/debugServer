package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func startServer() {

	router := mux.NewRouter()

	api := router.PathPrefix("/api/").Subrouter()
	api.Path("/add").HandlerFunc(HandleAddArticle).Methods("POST")
	api.Path("/clear").HandlerFunc(HandleClearArticles).Methods("POST")
	api.Path("/frontpage").HandlerFunc(HandleGetFrontpage).Methods("GET")
	api.Path("/article/{title}").HandlerFunc(HandleGetArticle).Methods("GET")

	path := "./src/github.com/opinionated/debugServer/debugFrontEnd"
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(path)))

	http.ListenAndServe(":8002", router)
}

// returns the handler for the api
func apiHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/add", HandleAddArticle).Methods("POST")
	router.HandleFunc("/clear", HandleClearArticles).Methods("POST")
	router.HandleFunc("/frontpage", HandleGetFrontpage).Methods("GET")
	router.HandleFunc("/article/{title}", HandleGetArticle).Methods("GET")
	return router
}
