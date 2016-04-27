package debugAPI_test

import (
	"encoding/json"
	"github.com/opinionated/analyzer-core/analyzer"
	"github.com/opinionated/articleStore"
	"github.com/opinionated/debugServer/debugAPI"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildStore(t *testing.T) {
	t.Skip("unskip to rebuild test data")
	store := articleStore.BuildStore(
		"./tmp",
		"json")

	createAndStore := func(title string, taxScore int) {
		body := title
		article := debugAPI.GenericArticle{
			Title: title,
			Body:  body,
		}

		store.CreateFolder(title)
		data, err := json.Marshal(article)
		assert.Nil(t, err)
		assert.Nil(t, store.StoreData(data, "Body", article.Title))

		// now store the fake taxonomy
		tax := map[string]int{
			"score": taxScore,
		}

		data, err = json.Marshal(tax)
		assert.Nil(t, err)
		assert.Nil(t, store.StoreData(data, "Taxonomy", article.Title))

	}

	createAndStore("a", 1)
	createAndStore("b", 2)
	createAndStore("c", 3)
	createAndStore("d", 2)
	createAndStore("e", 1)
}

func TestToDebug(t *testing.T) {
	store := articleStore.BuildStore(
		"./tmp",
		"json")
	debugAPI.SetStore(store)

	analyzed := analyzer.Analyzable{FileName: "a"}
	related := []analyzer.Analyzable{
		analyzer.Analyzable{FileName: "b"},
		analyzer.Analyzable{FileName: "c"},
	}

	article, err := debugAPI.ToDebug(analyzed, related)
	assert.Nil(t, err)

	assert.Equal(t, "a", article.Title)
	assert.Len(t, article.Related, 2)
}

func prepTest(t *testing.T) {
	debugAPI.SetServerURL("http://localhost:8002")
	store := articleStore.BuildStore(
		"./tmp",
		"json")
	debugAPI.SetStore(store)
	assert.Nil(t, debugAPI.Clear())
}

func TestPush(t *testing.T) {
	prepTest(t)

	analyzed := analyzer.Analyzable{FileName: "a"}
	related := []analyzer.Analyzable{
		analyzer.Analyzable{FileName: "b"},
		analyzer.Analyzable{FileName: "c"},
	}

	article, err := debugAPI.ToDebug(analyzed, related)
	assert.Nil(t, err)

	err = debugAPI.Push(article)
	assert.Nil(t, err)
}
