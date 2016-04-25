package main

import (
	"encoding/json"
	"github.com/opinionated/debugServer/debugAPI"
)

// for holding all the articles
var articles = articleList{
	limit:    10,
	count:    0,
	titleMap: make(map[string]debugAPI.GenericArticle),
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

	titleMap map[string]debugAPI.GenericArticle
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
func (list *articleList) push(article debugAPI.GenericArticle) {
	node := new(articleNode)
	node.article = article
	node.next = list.start

	list.start = node
	if list.end == nil {
		list.end = node
	}

	list.titleMap[article.Title] = article
	for _, related := range article.Related {
		list.titleMap[related.Title] = related
	}

	// do we need to bump
	if list.count == list.limit {
		list.popBack()
	} else {
		list.count++
	}
}

// get an article from it's title
func (list articleList) articleByTitle(title string) debugAPI.GenericArticle {
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
	article debugAPI.GenericArticle
}
