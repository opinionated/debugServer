package debugAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opinionated/analyzer-core/analyzer"
	"github.com/opinionated/articleStore"
	"net/http"
)

// GenericArticle stores all teh article stuff on the server
type GenericArticle struct {
	Title string `json:"Title"`
	// TODO: make sure that this decodes properly...
	Body  string `json:"Body"` // ignore this field
	Blurb string `json:"Blurb"`

	// articles related to the main article
	Related []GenericArticle `json:"Related,omitempty"`

	// for the first send json stuff
	DebugInfo map[string]interface{}
}

// ToDebug converts an analyzable article into a debug article
func ToDebug(analyzed analyzer.Analyzable, related []analyzer.Analyzable) (GenericArticle, error) {

	ret := GenericArticle{}

	ok, err := store.FolderExists(analyzed.FileName)
	if err != nil {
		return ret, fmt.Errorf("error checking for file: %s", err.Error())
	} else if !ok {
		return ret, fmt.Errorf("article %s does not exists", analyzed.FileName)
	}

	data, err := store.GetData("Body", analyzed.FileName)
	if err != nil {
		return ret, fmt.Errorf("error getting body: %s", err.Error())
	}

	// need to do this because of inconsistencies in the json
	var parsed map[string]string
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return ret, err
	}

	// actually build the article
	ret.Title = parsed["Title"]
	ret.Body = parsed["Data"]
	ret.Blurb = parsed["Desciption"]
	ret.Related = make([]GenericArticle, len(related))
	ret.DebugInfo = make(map[string]interface{})

	for i := range related {
		rel, err := ToDebug(related[i], nil)
		if err != nil {
			return ret,
				fmt.Errorf("trying to build %s, error on %s: %s",
					ret.Title, related[i].FileName, err.Error())
		}
		ret.Related[i] = rel
	}

	return ret, err
}

// Push an article onto the server
func Push(article GenericArticle) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}

	buff := bytes.NewBuffer(data)
	fmt.Println("the url is:", serverURL, "pushing:", article)

	resp, err := http.Post(serverURL+"/api/add", "application/json", buff)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return err
}

// Clear all the articles from the server
func Clear() error {
	_, err := http.Post(serverURL+"/api/clear", "", nil)
	return err
}

var store articleStore.Store
var serverURL string

// SetStore of article bodies
func SetStore(s articleStore.Store) {
	store = s
}

// SetServerURL give this the location of the server
func SetServerURL(url string) {
	serverURL = url
}
