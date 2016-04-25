package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func asbytes(vars ...interface{}) []byte {
	str := fmt.Sprint(vars...)
	return []byte(str)
}

// HandleAddArticle takes care of adding an article to article list.
// The article(s) to add are passed in as the JSON body.
// Endpoint should be /add.
func HandleAddArticle(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(asbytes("error parsing input:", err.Error()))
		return
	}

	var article genericArticle
	err = json.Unmarshal(raw, &article)
	if err != nil {
		w.Write(asbytes("error parsing body:", err.Error()))
		return
	}

	articles.push(article)
}

// HandleGetFrontpage returns a list of all the "top" articles.
// TODO: decide if this is how we want to do this perminently.
// Endpoint should be /frontpage
func HandleGetFrontpage(w http.ResponseWriter, r *http.Request) {
	data, err := articles.buildJSON()
	if err != nil {
		w.Write(asbytes("error converting to json:", err.Error()))
		return
	}

	n, err := w.Write(data)
	if err != nil {
		w.Write(asbytes("error writing data:", err.Error()))
	} else if n != len(data) {
		w.Write(asbytes("error writing data: did not write full json"))
	}
}

// HandleGetArticle returns an article body and list of related articles.
// It should be called after the article is clicked on.
// The endpoint should be /article/{title}
func HandleGetArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	article, ok := articles.titleMap[vars["title"]]
	if !ok {
		w.Write(asbytes("error finding article called", vars["title"]))
		return
	}

	ret := make(map[string]interface{})
	ret["Body"] = article.Body

	related := make([]map[string]string, len(article.Related))
	for i, r := range article.Related {
		related[i] = map[string]string{
			"Title": r.Title,
		}
	}

	ret["Related"] = related

	data, err := json.Marshal(ret)
	if err != nil {
		w.Write(asbytes("error marshalling json"))
		return
	}

	w.Write(data)
}
