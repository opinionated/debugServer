package main

import (
	"bytes"
	"encoding/json"
	"github.com/opinionated/debugServer/debugAPI"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func getHTTPResponse(url string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	newHandler().ServeHTTP(recorder, req)

	return recorder
}

func postHTTPResponse(url, body string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	newHandler().ServeHTTP(recorder, req)

	return recorder
}

// makes a simple test article
// store everything in a json map to make it easier
func makeArticle(title string, body string,
	related []map[string]interface{}) map[string]interface{} {

	article := make(map[string]interface{})
	article["Title"] = title
	article["Body"] = body
	if related != nil {
		article["Related"] = related
	}
	return article
}

// encodes to json then converts to string
func setToString(set map[string]interface{}) (string, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(set)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

var testSet = struct {
	guns map[string]interface{}
	god  map[string]interface{}
}{
	guns: makeArticle("Guns", "guns r great",
		[]map[string]interface{}{
			makeArticle("Guns good", "shoot stuff", nil),
			makeArticle("Stop guns", "stop em", nil),
		}),
	god: makeArticle("God", "separate church and state",
		[]map[string]interface{}{
			makeArticle("Abortions bad", "god said so!", nil),
			makeArticle("America christian", "one nation under god", nil),
		}),
}

func addTestSet(t *testing.T) {
	clearArticles()
	// add sets
	gunstr, err := setToString(testSet.guns)
	assert.Nil(t, err)
	resp := postHTTPResponse("/add", gunstr)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "", resp.Body.String())

	godstr, err := setToString(testSet.god)
	assert.Nil(t, err)
	resp = postHTTPResponse("/add", godstr)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "", resp.Body.String())

	assert.Equal(t, 2, articles.count)
}

func clearArticles() {
	articles = articleList{
		limit:    10,
		count:    0,
		titleMap: make(map[string]debugAPI.GenericArticle),
		start:    nil,
		end:      nil,
	}
}

func TestAddArticle(t *testing.T) {
	clearArticles()
	str, err := setToString(testSet.guns)
	assert.Nil(t, err)

	resp := postHTTPResponse("/add", str)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "", resp.Body.String())

	// now go test the body
	assert.Equal(t, 1, articles.count)
	assert.NotNil(t, articles.start)

	article := articles.start.article
	assert.Equal(t, article.Title, "Guns")
	assert.Len(t, article.Related, 2)
	assert.Equal(t, article.Body, testSet.guns["Body"])
}

func TestAddToPop(t *testing.T) {
	clearArticles()
	str, err := setToString(testSet.guns)
	assert.Nil(t, err)

	for i := 0; i < articles.limit+1; i++ {
		resp := postHTTPResponse("/add", str)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "", resp.Body.String())
	}
}

func TestGetArticles(t *testing.T) {
	addTestSet(t)

	resp := getHTTPResponse("/frontpage")
	assert.Equal(t, 200, resp.Code)
	assert.NotNil(t, resp.Body)

	var articles []debugAPI.GenericArticle
	err := json.Unmarshal(resp.Body.Bytes(), &articles)

	assert.Nil(t, err)
	assert.Len(t, articles, 2)
	assert.Equal(t, "God", articles[0].Title)
	assert.Equal(t, "Guns", articles[1].Title)

	assert.Nil(t, articles[0].Related)
	assert.Nil(t, articles[1].Related)
}

func TestGetArticle(t *testing.T) {
	addTestSet(t)

	resp := getHTTPResponse("/article/Guns")
	assert.Equal(t, 200, resp.Code)
	assert.NotNil(t, resp.Body)

	var article debugAPI.GenericArticle
	err := json.Unmarshal(resp.Body.Bytes(), &article)
	assert.Nil(t, err)

	assert.Equal(t, testSet.guns["Body"], article.Body)
	assert.Len(t, article.Related, 2)
	assert.Equal(t, "Guns good", article.Related[0].Title)
	assert.Equal(t, "Stop guns", article.Related[1].Title)
}

func TestGetRelated(t *testing.T) {
	addTestSet(t)

	resp := getHTTPResponse("/article/Guns good")
	assert.Equal(t, 200, resp.Code)
	assert.NotNil(t, resp.Body)

	if _, ok := articles.titleMap["Guns good"]; !ok {
		t.Error("Expected article, but could not find!")
		t.Error("list is:", articles.titleMap)
	}

	var article debugAPI.GenericArticle
	err := json.Unmarshal(resp.Body.Bytes(), &article)
	if !assert.Nil(t, err) {
		t.Error("json body is:", resp.Body.String())
	}

	assert.Equal(t, "shoot stuff", article.Body)
	assert.Len(t, article.Related, 0)
}

func TestAddDebug(t *testing.T) {
	guns := testSet.guns
	guns["Title"] = "debug"
	guns["DebugInfo"] = map[string]interface{}{
		"Score": 1.0,
		"Taxonomy": map[string]interface{}{
			"a": 1,
			"b": 2,
		},
	}

	str, err := setToString(guns)
	assert.Nil(t, err)

	// add our article with debug info up to the server
	clearArticles()
	resp := postHTTPResponse("/add", str)
	assert.NotNil(t, resp.Body)
	assert.Equal(t, "", resp.Body.String())

	resp = getHTTPResponse("/article/debug")
	assert.NotNil(t, resp.Body)

	var data map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &data)
	assert.Nil(t, err)

	// pull the debug info out of the json
	rawDebug, ok := data["DebugInfo"]
	assert.True(t, ok)
	debugInfo, ok := rawDebug.(map[string]interface{})
	assert.True(t, ok)

	// check that the sub info exists
	assert.Equal(t, 1.0, debugInfo["Score"])
	taxonomy, ok := debugInfo["Taxonomy"].(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, 1.0, taxonomy["a"])
	assert.Equal(t, 2.0, taxonomy["b"])
}
