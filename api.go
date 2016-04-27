package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/opinionated/debugServer/debugAPI"
	"io/ioutil"
	"net/http"
)

// HandleAddArticle takes care of adding an article to article list.
// The article(s) to add are passed in as the JSON body.
// Endpoint should be api/add.
func HandleAddArticle(w http.ResponseWriter, r *http.Request) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(asbytes("error parsing input:", err.Error()))
		return
	}

	// parse and push onto the cache
	var article debugAPI.GenericArticle
	err = json.Unmarshal(raw, &article)
	if err != nil {
		w.Write(asbytes("error parsing body:", err.Error()))
		return
	}

	cache.lock()
	cache.push(article)
	cache.unlock()
}

// HandleGetFrontpage returns a list of all the "top" articles.
// TODO: decide if this is how we want to do this perminently.
// The body will look something like this:
// [ {Title:"a"}, {Title:"b"} ]
// Endpoint should be api/frontpage.
func HandleGetFrontpage(w http.ResponseWriter, r *http.Request) {

	// START CRITICAL SECTION
	cache.lock()

	// convert the list of articles to JSON
	articleMap := make([]map[string]string, cache.count)
	tmp := cache.start
	i := 0
	for tmp != nil {
		articleMap[i] = map[string]string{
			"Title": tmp.article.Title,
		}

		i++
		tmp = tmp.next
	}

	cache.unlock()
	// END CRITICAL SECTION

	// marshal and send the data
	data, err := json.Marshal(articleMap)
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

	if err != nil {
		w.Write(asbytes("error setting header!"))
	}
}

// HandleGetArticle returns an article body and list of related articles.
// It should be called after the article is clicked on.
// The body will look something like this:
// { Body:"...", DebugInfo:{}, Related:[{Title:"", DebugInfo:""}]}
// note that debug info can have anything in it and related is an array
// The endpoint should be api/article/{title}
func HandleGetArticle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	cache.lock()
	article, ok := cache.articleByTitle(vars["title"])
	cache.unlock()

	if !ok {
		w.Write(asbytes("error finding article called", vars["title"]))
		return
	}

	// doesn't yet have the body or the debug info
	ret := make(map[string]interface{})
	ret["Body"] = article.Body
	ret["DebugInfo"] = article.DebugInfo
	ret["Title"] = article.Title

	// build the related articles only if they exist
	if len(article.Related) > 0 {

		related := make([]map[string]interface{}, len(article.Related))
		for i, r := range article.Related {

			// send the debug info with the article
			related[i] = map[string]interface{}{
				"Title":     r.Title,
				"DebugInfo": r.DebugInfo,
			}
		}

		ret["Related"] = related
	}

	data, err := json.Marshal(ret)
	if err != nil {
		w.Write(asbytes("error marshalling json"))
		return
	}

	w.Write(data)
}

// HandleClearArticles dumps the cache.
// By default the cache holds 10 articles on the home page.
// The endpoint should be api/clear
func HandleClearArticles(w http.ResponseWriter, r *http.Request) {
	cache.lock()
	cache.clear()
	cache.unlock()
}

func asbytes(vars ...interface{}) []byte {
	// helper converts a string to bytes for writing msgs
	str := fmt.Sprint(vars...)
	return []byte(str)
}
