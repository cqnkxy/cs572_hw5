package views

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_solr "github.com/rtt/Go-Solr"
)

func Suggest(w http.ResponseWriter, r *http.Request) {
	words, ok := r.URL.Query()["words"]
	log.Println("Got suggest request for ", words)
	if !ok {
		http.NotFound(w, r)
	}
	bytes, err := _solr.HTTPGet(fmt.Sprintf(
		"http://localhost:8983/solr/newcore/suggest?indent=on&q=%s&wt=json",
		words[0],
	))
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
		js, _ := json.Marshal([]string{})
		fmt.Fprint(w, js)
		return
	}

	var container interface{}
	if err := json.Unmarshal(bytes, &container); err != nil {
		log.Println(err)
		js, _ := json.Marshal([]string{})
		fmt.Fprint(w, js)
		return
	}

	response_root := container.(map[string]interface{})
	response := response_root["suggest"].(map[string]interface{})["suggest"].(map[string]interface{})[words[0]].(map[string]interface{})["suggestions"].([]interface{})
	suggestions := []string{}
	for i := 0; i < len(response); i += 1 {
		term := response[i].(map[string]interface{})["term"].(string)
		if len(term) < 15 {
			suggestions = append(suggestions, term)
		}
	}
	if js, err := json.Marshal(suggestions); err == nil {
		log.Println("Suggest as ", suggestions)
		w.Write(js)
	} else {
		log.Println(err)
	}
}
