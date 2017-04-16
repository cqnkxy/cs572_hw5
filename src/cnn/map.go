package cnn

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var IdToUrl map[string]string

func init() {
	IdToUrl = map[string]string{}
	file, err := os.Open("/Users/xueyuan/Documents/USC/csci572/hw5/src/data/mapCNNDataFile.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splitted := strings.Split(scanner.Text(), ",")
		IdToUrl[splitted[0]] = splitted[1]
	}
}
