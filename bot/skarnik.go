package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func translate(searchTerm string) (*string, error) {
	cleanSearchTerm := cleanTerm(searchTerm)
	firstWord := strings.Fields(cleanSearchTerm)[0]

	resp, err := requestSkarnik(firstWord)
	if err != nil {
		return nil, err
	}

	translation, err := parseSkarnikResponse(resp)
	if err != nil {
		return nil, err
	}

	return translation, nil
}

func parseSkarnikResponse(resp *http.Response) (*string, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
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

	return &translation, nil
}

func requestSkarnik(searchTerm string) (*http.Response, error) {
	requestUrl := fmt.Sprintf("https://www.skarnik.by/search_json?term=%s&lang=rus", searchTerm)

	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var suggestions []Suggestion
	err = json.Unmarshal(body, &suggestions)
	if err != nil {
		return nil, err
	}

	if len(suggestions) == 0 {
		return nil, errors.New("No results found")
	}

	requestUrl = fmt.Sprintf("https://www.skarnik.by/rusbel/%d", suggestions[0].Id)

	resp, err = http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
