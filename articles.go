package main

import (
	"github.com/opinionated/debugServer/debugAPI"
	"sync"
)

// for holding all the articles
var cache = NewCache()

// ArticleCache stores all the active articles
// the caller must protect the cache
type ArticleCache struct {
	start *articleNode
	end   *articleNode

	// how many can be in the list
	limit int
	count int

	titleMap map[string]articlePair
	mutex    *sync.Mutex
}

// NewCache creates a new cache
func NewCache() ArticleCache {
	return ArticleCache{
		limit:    10,
		count:    0,
		titleMap: make(map[string]articlePair),
		start:    nil,
		end:      nil,
		mutex:    new(sync.Mutex),
	}
}

// note that these are non-recursive
func (list *ArticleCache) lock() {
	list.mutex.Lock()
}

func (list *ArticleCache) unlock() {
	list.mutex.Unlock()
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
		list.pop()
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

func (list *ArticleCache) clear() {
	list.start = nil
	list.end = nil
	list.count = 0

	// this is safe
	for key := range list.titleMap {
		delete(list.titleMap, key)
	}
}

// get an article by it's title
func (list ArticleCache) articleByTitle(title string) (debugAPI.GenericArticle, bool) {
	pair, ok := list.titleMap[title]
	return pair.article, ok
}

func (list *ArticleCache) pop() {
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
	list.count--
}

// forward list
type articleNode struct {
	next    *articleNode
	article debugAPI.GenericArticle
}
