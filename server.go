package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func startServer(port string) {

	router := mux.NewRouter()

	api := router.PathPrefix("/api/").Subrouter()
	api.Path("/add").HandlerFunc(HandleAddArticle).Methods("POST")
	api.Path("/clear").HandlerFunc(HandleClearArticles).Methods("POST")
	api.Path("/frontpage").HandlerFunc(HandleGetFrontpage).Methods("GET")
	api.Path("/article/{title}").HandlerFunc(HandleGetArticle).Methods("GET")

	// otherwise go to the file server
	router.PathPrefix("/{filename}").HandlerFunc(HandleServeFile).Methods("GET")

	http.ListenAndServe(port, router)
}

// HandleServeFile takes care of appending html onto certain articles
func HandleServeFile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	filename, ok := vars["filename"]

	if !ok {
		w.Write(asbytes("oh nose! bad bad bad!"))
		return
	}

	path := "./src/github.com/opinionated/debugServer/debugFrontEnd/"
	if strings.Contains(filename, ".") {
		http.ServeFile(w, r, path+filename)
	} else {
		http.ServeFile(w, r, path+filename+".html")
	}
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
