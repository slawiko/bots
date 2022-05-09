package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const EmptyResultMessage = "Нічога не знайшоў :("

type Suggestion struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}

func cleanTerm(searchTerm string) string {
	cleanSearchTerm := strings.ReplaceAll(searchTerm, "ў", "щ")
	cleanSearchTerm = strings.ReplaceAll(cleanSearchTerm, "і", "и")
	cleanSearchTerm = strings.ReplaceAll(cleanSearchTerm, "’", "ъ")
	cleanSearchTerm = strings.ReplaceAll(cleanSearchTerm, "'", "ъ")

	return cleanSearchTerm
}

func translate(searchTerm string) string {
	cleanSearchTerm := cleanTerm(searchTerm)

	firstWord := strings.Fields(cleanSearchTerm)[0]

	requestUrl := fmt.Sprintf("https://www.skarnik.by/search_json?term=%s&lang=rus", firstWord)

	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var suggestions []Suggestion
	json.Unmarshal(body, &suggestions)

	if len(suggestions) == 0 {
		return EmptyResultMessage
	}

	requestUrl = fmt.Sprintf("https://www.skarnik.by/rusbel/%d", suggestions[0].Id)

	resp, err = http.Get(requestUrl)
	if err != nil {
		log.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
	}

	section := doc.Find("#trn")

	translation := ""

	section.Find("font[color=\"831b03\"]").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			translation += s.Text()
		} else {
			translation += fmt.Sprintf(", %s", s.Text())
		}
	})

	log.Println(translation)

	return translation
}
