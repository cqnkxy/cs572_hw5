package views

import (
	"encoding/json"
	"log"
	"testing"

	_solr "github.com/rtt/Go-Solr"
)

func TestSuggest(t *testing.T) {
	log.Println("Testing suggest...")
	bytes, err := _solr.HTTPGet("http://localhost:8983/solr/newcore/suggest?indent=on&q=ca&wt=json")
	if err != nil {
		panic(err)
	}

	var container interface{}
	if err := json.Unmarshal(bytes, &container); err != nil {
		panic(err)
	}

	response_root := container.(map[string]interface{})
	response := response_root["suggest"].(map[string]interface{})["suggest"].(map[string]interface{})["ca"].(map[string]interface{})["suggestions"].([]interface{})
	log.Println(response)
	suggestions := []string{}
	for i := 0; i < len(response); i += 1 {
		term := response[i].(map[string]interface{})["term"].(string)
		suggestions = append(suggestions, term)
	}
	log.Println("Got suggestions for ca as ", suggestions)
	log.Println("Passed")
}
