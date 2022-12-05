package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Suggestion struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

func Translate(searchTerm string, isDetailed bool) (string, error) {
	words := strings.Fields(searchTerm)

	suggestions, err := getSkarnikSuggestions(words[0])
	if err != nil {
		log.Println(err)
		return "", err
	}
	if len(suggestions) == 0 {
		return "", errors.New("no translation found")
	}

	resp, err := requestSkarnik(suggestions[0])
	if err != nil {
		log.Println(err)
		return "", err
	}

	if isDetailed {
		return DetailedTranslationParse(resp.Body)
	} else {
		_, htmltranslation, err := ShortTranslationParse(resp.Body)
		return htmltranslation, err
	}
}

func getSkarnikSuggestions(searchTerm string) ([]Suggestion, error) {
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

	return suggestions, nil
}

func requestSkarnik(suggestion Suggestion) (*http.Response, error) {
	requestUrl := fmt.Sprintf("https://www.skarnik.by/rusbel/%d", suggestion.ID)

	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
