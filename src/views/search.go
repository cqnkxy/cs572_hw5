package views

import (
	"encoding/json"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"cnn"

	_solr "github.com/rtt/Go-Solr"
)

const (
	CNN_DATA_PATH = "/Users/xueyuan/Documents/USC/csci572/hw4/CNNData/CNNDownloadData/"
)

var skipRe *regexp.Regexp

func init() {
	skipRe = regexp.MustCompile("^script|style|<img|head|title")
}

type SearchResult struct {
	Title   string
	Url     string
	Id      string
	Snippet string
}

func traverse(n *html.Node, queryRe *regexp.Regexp) (string, bool) {
	if skipRe.MatchString(n.Data) {
		return "", false
	}
	if n.Type == html.TextNode && queryRe.MatchString(strings.ToLower(n.Data)) {
		return n.Data, true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result, ok := traverse(c, queryRe); ok {
			return result, true
		}
	}
	return "", false
}

func findFisrtMatchingSentence(query []string, fileId string) string {
	file, err := os.Open(path.Join(CNN_DATA_PATH, fileId))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	doc, err := html.Parse(file)
	if err != nil {
		panic(err)
	}
	queryRe := regexp.MustCompile("(?i)" + strings.Join(query, " "))
	if res, ok := traverse(doc, queryRe); ok {
		return res
	}
	return "N.A."
}

func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()["query"]
	log.Println("Get request for ", query,
		" using ", r.URL.Query()["method"])
	q := _solr.Query{
		Params: _solr.URLParamMap{
			"q":      query,
			"indent": []string{"on"},
			"wt":     []string{"json"},
		},
		Rows: 10,
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
		title, id, description, url := "", "N.A.", "N.A.", "N.A."
		ids := strings.Split(results.Get(i).Field("id").(string), "/")
		id = ids[len(ids)-1]
		if titles, ok := results.Get(i).Field("title").([]interface{}); ok && len(titles) > 0 {
			title = titles[0].(string)
		}
		if descriptions, ok := results.Get(i).Field("description").([]interface{}); ok && len(descriptions) > 0 {
			description = descriptions[0].(string)
			if !strings.Contains(description, strings.Join(query, " ")) {
				description = findFisrtMatchingSentence(query, id)
			}
		}
		url = cnn.IdToUrl[id]
		searchResults = append(searchResults, &SearchResult{
			Url:     url,
			Title:   title,
			Snippet: description,
			Id:      id,
		})
	}
	if js, err := json.Marshal(searchResults); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
