package main

import (
	"encoding/json"
)

// for holding all the articles
var articles = articleList{
	limit:    10,
	count:    0,
	titleMap: make(map[string]genericArticle),
	start:    nil,
	end:      nil,
}

// holds the list of active articles
// not currently thread safe...
type articleList struct {
	start *articleNode
	end   *articleNode

	// how many can be in the list
	limit int
	count int

	titleMap map[string]genericArticle
}

// buildJSON converts the list of articles to JSON
func (list *articleList) buildJSON() ([]byte, error) {
	articles := make([]map[string]string, list.count)
	tmp := list.start
	i := 0
	for tmp != nil {
		articles[i] = map[string]string{
			"Title": tmp.article.Title,
		}

		i++
		tmp = tmp.next
	}

	return json.Marshal(articles)
}

// add a new article to the list
func (list *articleList) push(article genericArticle) {
	node := new(articleNode)
	node.article = article
	node.next = list.start

	list.start = node
	if list.end == nil {
		list.end = node
	}

	list.titleMap[article.Title] = article

	// do we need to bump
	if list.count == list.limit {
		list.popBack()
	} else {
		list.count++
	}
}

// get an article from it's title
func (list articleList) articleByTitle(title string) genericArticle {
	return list.titleMap[title]
}

func (list *articleList) popBack() {
	tmp := list.start
	for tmp.next != list.end {
		tmp = tmp.next
	}

	// remove from map
	toRemove := tmp.next
	delete(list.titleMap, toRemove.article.Title)
	related := toRemove.article.Related
	for i := range related {
		delete(list.titleMap, related[i].Title)
	}

	// bump from list
	list.end = tmp
	list.end.next = nil

}

// forward list
type articleNode struct {
	next    *articleNode
	article genericArticle
}

// generic article for all types
type genericArticle struct {
	Title string `json:"Title"`
	// TODO: make sure that this decodes properly...
	Body  string `json:"Body"` // ignore this field
	Blurb string `json:"Blurb"`

	// articles related to the main article
	Related []genericArticle `json:"Related,omitempty"`

	// for the first send json stuff
	shortMap map[string]string
}
