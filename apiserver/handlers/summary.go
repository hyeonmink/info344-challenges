package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"encoding/json"

	"errors"

	"golang.org/x/net/html"
)

//openGraphPrefix is the prefix used for Open Graph meta properties
const openGraphPrefix = "og:"

//openGraphProps represents a map of open graph property names and values
type openGraphProps map[string]string

func getPageSummary(url string) (openGraphProps, error) {
	//Get the URL
	//If there was an error, return it
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v\n", err)
	}

	//ensure that the response body stream is closed eventually
	//HINTS: https://gobyexample.com/defer
	//https://golang.org/pkg/net/http/#Response
	defer resp.Body.Close()

	//if the response StatusCode is >= 400
	//return an error, using the response's .Status
	//property as the error message
	if resp.StatusCode >= 400 {
		err = errors.New(resp.Status)
		return nil, err
	}

	//if the response's Content-Type header does not
	//start with "text/html", return an error noting
	//what the content type was and that you were
	//expecting HTML
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, fmt.Errorf("response content type was %s not text/html\n", ctype)
	}

	//create a new openGraphProps map instance to hold
	//the Open Graph properties you find
	//(see type definition above)
	m := make(openGraphProps)

	//tokenize the response body's HTML and extract
	//any Open Graph properties you find into the map,
	//using the Open Graph property name as the key, and the
	//corresponding content as the value.
	//strip the openGraphPrefix from the property name before
	//you add it as a new key, so that the key is just `title`
	//and not `og:title` (for example).
	tokenizer := html.NewTokenizer(resp.Body)
	err = nil
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err = tokenizer.Err()
			return m, err
		} else if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()

			if "meta" == token.Data {
				attr := token.Attr
				key := ""
				val := ""

				for i := 0; i < len(attr); i++ {
					temp := attr[i]
					if "property" == temp.Key {
						key = strings.TrimPrefix(temp.Val, openGraphPrefix)
					} else if "content" == temp.Key {
						val = temp.Val
					} else if "name" == temp.Key && m["description"] == "" {
						key = "description"
					}
					if val != "" && key != "" {
						m[key] = val
					}
				}
			} else if "title" == token.Data && m["title"] == "" {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					token := tokenizer.Token()
					m["title"] = token.Data
				}
			}
		} else if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if "head" == token.Data {
				return m, err
			}
		}

	}
}

//HINTS: https://info344-s17.github.io/tutorials/tokenizing/
//https://godoc.org/golang.org/x/net/html

//SummaryHandler fetches the URL in the `url` query string parameter, extracts
//summary information about the returned page and sends those summary properties
//to the client as a JSON-encoded object.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	//Add the following header to the response
	//   Access-Control-Allow-Origin: *
	//this will allow JavaScript served from other origins
	//to call this API
	// w.Header().Add("Content-Type", "application/json; charset=utf-8")
	// w.Header().Add("Access-Control-Allow-Origin", "*")

	//get the `url` query string parameter
	//if you use r.FormValue() it will also handle cases where
	//the client did POST with `url` as a form field
	//HINT: https://golang.org/pkg/net/http/#Request.FormValue
	url := r.URL.Query().Get("url")

	//if no `url` parameter was provided, respond with
	//an http.StatusBadRequest error and return
	//HINT: https://golang.org/pkg/net/http/#Error
	if len(url) == 0 {
		http.Error(w, "no url", http.StatusBadRequest)
		return
	}

	//call getPageSummary() passing the requested URL
	//and holding on to the returned openGraphProps map
	//(see type definition above)
	m, err := getPageSummary(url)

	//if you get back an error, respond to the client
	//with that error and an http.StatusBadRequest code
	if err != nil {
		http.Error(w, "error fetching URL: "+err.Error(), http.StatusBadRequest)
	}

	//otherwise, respond by writing the openGrahProps
	//map as a JSON-encoded object
	//add the following headers to the response before
	//you write the JSON-encoded object:
	//   Content-Type: application/json; charset=utf-8
	//this tells the client that you are sending it JSON
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(m); err != nil {
		http.Error(w, "error encoding json: "+err.Error(), http.StatusInternalServerError)
	}
}
