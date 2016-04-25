package main

import (
	"github.com/opinionated/debugServer/debugAPI"
)

// for holding all the articles
var cache = ArticleCache{
	limit:    10,
	count:    0,
	titleMap: make(map[string]articlePair),
	start:    nil,
	end:      nil,
}

// ArticleCache stores all the active articles
type ArticleCache struct {
	start *articleNode
	end   *articleNode

	// how many can be in the list
	limit int
	count int

	titleMap map[string]articlePair
}

// because related articles may appear multiple times
type articlePair struct {
	count   int
	article debugAPI.GenericArticle
}

// add a new article to the list
func (list *ArticleCache) push(article debugAPI.GenericArticle) {
	node := new(articleNode)
	node.article = article
	node.next = list.start

	list.start = node
	if list.end == nil {
		list.end = node
	}

	list.addToMap(article)

	// do we need to bump
	if list.count == list.limit {
		list.popBack()
	} else {
		list.count++
	}
}

func (list *ArticleCache) addToMap(article debugAPI.GenericArticle) {
	list.titleMap[article.Title] = articlePair{article: article, count: 1}

	for _, related := range article.Related {
		pair, ok := list.titleMap[related.Title]
		if !ok {
			list.titleMap[related.Title] =
				articlePair{article: related, count: 1}
		} else {
			pair.count++
			list.titleMap[related.Title] = pair
		}
	}
}

func (list *ArticleCache) removeFromMap(article debugAPI.GenericArticle) {
	pair := list.titleMap[article.Title]
	if pair.count == 1 {
		delete(list.titleMap, article.Title)
	} else {
		pair.count--
		list.titleMap[article.Title] = pair
	}
}

// get an article from it's title
func (list ArticleCache) articleByTitle(title string) (debugAPI.GenericArticle, bool) {
	pair, ok := list.titleMap[title]
	return pair.article, ok
}

func (list *ArticleCache) popBack() {
	tmp := list.start
	for tmp.next != list.end {
		tmp = tmp.next
	}

	// remove from map
	toRemove := tmp.next
	related := toRemove.article.Related
	for i := range related {
		list.removeFromMap(related[i])
	}
	list.removeFromMap(toRemove.article)

	// bump from list
	list.end = tmp
	list.end.next = nil

}

// forward list
type articleNode struct {
	next    *articleNode
	article debugAPI.GenericArticle
}
