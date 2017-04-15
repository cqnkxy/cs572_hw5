package views

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"spell"
)

func Correct(w http.ResponseWriter, r *http.Request) {
	words, ok := r.URL.Query()["words"]
	if !ok {
		http.NotFound(w, r)
	}
	corrected := []string{}
	for _, word := range strings.Split(words[0], " ") {
		if word != "" {
			corrected = append(corrected, spell.Correct(word))
		}
	}
	log.Printf("Got correct request for %s, correted to %s\n", words[0], strings.Join(corrected, " "))
	fmt.Fprint(w, strings.Join(corrected, " "))
}
