package spell

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var model map[string]int

func init() {
	// model = train("/Users/xueyuan/Documents/USC/csci572/hw5/src/data/big.txt")
}

func train(training_data string) map[string]int {
	log.Println("Training data...")
	NWORDS := make(map[string]int)
	pattern := regexp.MustCompile("[a-z]+")
	if content, err := ioutil.ReadFile(training_data); err == nil {
		for _, w := range pattern.FindAllString(strings.ToLower(string(content)), -1) {
			NWORDS[w]++
		}
	} else {
		panic(err)
	}
	log.Println("Training finished!")
	return NWORDS
}

// func train(training_data string) map[string]int {
// 	log.Println("Training data...")
// 	file, err := os.Open(training_data)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()
// 	NWORDS := make(map[string]int)
// 	scanner := bufio.NewScanner(file)
// 	pattern := regexp.MustCompile("[a-z]+")
// 	for scanner.Scan() {
// 		for _, w := range pattern.FindAllString(strings.ToLower(scanner.Text()), -1) {
// 			NWORDS[w]++
// 		}
// 	}
// 	log.Println("Training finished!")
// 	return NWORDS
// }

func edits1(word string, ch chan string) {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	type Pair struct{ a, b string }
	var splits []Pair
	for i := 0; i < len(word)+1; i++ {
		splits = append(splits, Pair{word[:i], word[i:]})
	}

	for _, s := range splits {
		if len(s.b) > 0 {
			ch <- s.a + s.b[1:]
		}
		if len(s.b) > 1 {
			ch <- s.a + string(s.b[1]) + string(s.b[0]) + s.b[2:]
		}
		for _, c := range alphabet {
			if len(s.b) > 0 {
				ch <- s.a + string(c) + s.b[1:]
			}
		}
		for _, c := range alphabet {
			ch <- s.a + string(c) + s.b
		}
	}
}

func edits2(word string, ch chan string) {
	ch1 := make(chan string, 1024*1024)
	go func() { edits1(word, ch1); ch1 <- "" }()
	for e1 := range ch1 {
		if e1 == "" {
			break
		}
		edits1(e1, ch)
	}
}

func best(word string, edits func(string, chan string), model map[string]int) string {
	ch := make(chan string, 1024*1024)
	go func() { edits(word, ch); ch <- "" }()
	maxFreq := 0
	correction := ""
	for word := range ch {
		if word == "" {
			break
		}
		if freq, present := model[word]; present && freq > maxFreq {
			maxFreq, correction = freq, word
		}
	}
	return correction
}

func Correct(word string) string {
	if _, present := model[word]; present {
		return word
	}
	if correction := best(word, edits1, model); correction != "" {
		return correction
	}
	if correction := best(word, edits2, model); correction != "" {
		return correction
	}
	return word
}
