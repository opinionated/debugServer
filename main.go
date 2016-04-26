package main

import (
	"github.com/opinionated/debugServer/debugAPI"
)

func main() {
	// build the main articles
	guns := buildDebug("guns", "guns are great", []debugAPI.GenericArticle{
		buildDebug("guns kill", "look how many kids die from guns", nil),
		buildDebug("2nd amendment", "jesus said give us the guns", nil),
		buildDebug("bear attacks", "good thing they had a gun", nil),
		buildDebug("obama taking your guns", "dirty muslim", nil),
	})

	dolphins := buildDebug("5 most evil critters", "1-5:dolphins",
		[]debugAPI.GenericArticle{
			buildDebug("dolphins are nice", "said hostage of dolphins", nil),
			buildDebug("artic dolphin hunters", "most dangerous game", nil),
			buildDebug("tuna is mostly dolphin", "got heem", nil),
			buildDebug("dolphin attacks swimmers", "worse than shark", nil),
		})

	cache.push(dolphins)
	cache.push(guns)

	// spin up the server
	startServer()
}

func buildDebug(name string, body string,
	related []debugAPI.GenericArticle) debugAPI.GenericArticle {

	return debugAPI.GenericArticle{
		Title:   name,
		Body:    body,
		Related: related,
	}
}
