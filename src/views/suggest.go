package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	bs "bytes"

	_solr "github.com/rtt/Go-Solr"
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	words, ok := r.URL.Query()["words"]
	log.Println("Got suggest request for ", words)
	if !ok {
		http.Error(w, "Request format error!", http.StatusBadRequest)
		return
	}
	bytes, err := _solr.HTTPGet(fmt.Sprintf(
		"http://localhost:8983/solr/newcore/suggest?indent=on&q=%s&wt=json",
		words[0],
	))
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data struct {
		Suggest struct {
			Suggest map[string]struct{
				Suggestions []struct{
					Term string
				}
			}
		}
	}

	suggestions := []string{}
	if err = json.NewDecoder(bs.NewReader(bytes)).Decode(&data); err != nil {
		panic(err)
	}
	for _, suggestion := range data.Suggest.Suggest[words[0]].Suggestions {
		suggestions = append(suggestions, suggestion.Term)
	}
	if js, err := json.Marshal(suggestions); err == nil {
		log.Println("Suggest as ", suggestions)
		w.Write(js)
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
