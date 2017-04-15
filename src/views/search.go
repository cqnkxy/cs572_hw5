package views

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"cnn"

	_solr "github.com/rtt/Go-Solr"
)

type SearchResult struct {
	Title       string
	Url         string
	Description string
	Id          string
}

func Search(w http.ResponseWriter, r *http.Request) {
	log.Println("Get request for ", r.URL.Query()["query"],
		" using ", r.URL.Query()["method"])
	q := _solr.Query{
		Params: _solr.URLParamMap{
			"q":      r.URL.Query()["query"],
			"indent": []string{"on"},
			"wt":     []string{"json"},
		},
		Rows: 100,
	}
	if r.URL.Query()["method"][0] == "pagerank" {
		q.Sort = url.QueryEscape("pageRankFile desc")
	}
	res, err := solr.Select(&q)
	if err != nil {
		log.Println(err)
	}

	searchResults := []*SearchResult{}
	results := res.Results
	for i := 0; i < results.Len(); i += 1 {
		title, id, description, url := "N.A.", "N.A.", "N.A.", "N.A."
		ids := strings.Split(results.Get(i).Field("id").(string), "/")
		id = ids[len(ids)-1]
		if titles, ok := results.Get(i).Field("title").([]interface{}); ok && len(titles) > 0 {
			title = titles[0].(string)
		}
		if descriptions, ok := results.Get(i).Field("description").([]interface{}); ok && len(descriptions) > 0 {
			description = descriptions[0].(string)
		}
		url = cnn.IdToUrl[id]
		searchResults = append(searchResults, &SearchResult{
			Url:         url,
			Title:       title,
			Description: description,
			Id:          id,
		})
	}
	if js, err := json.Marshal(searchResults); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
