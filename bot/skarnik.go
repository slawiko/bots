package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
)

type Suggestion struct {
	ID    int    `json:"id"`
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
	words := strings.Fields(cleanSearchTerm)

	translation := ""

	for index, word := range words {
		suggestions, err := getScarnikSuggestions(word)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(suggestions) == 0 {
			continue
		}

		resp, err := requestSkarnik(suggestions[0])
		if err != nil {
			log.Println(err)
			continue
		}

		wordShortTranslation, err := parseTranslation(resp.Body)
		if err != nil {
			return nil, err
		}

		if index == 0 {
			translation = *wordShortTranslation
		} else {
			translation = translation + ";\n" + *wordShortTranslation
		}
	}

	if len(translation) == 0 {
		return nil, errors.New("No translation found")
	}

	return &translation, nil
}

func parseTranslation(body io.Reader) (detailedTranslation *string, err error) {
	tknzr := html.NewTokenizer(body)

	stack := stack{
		stack: make([]string, 0),
	}

	detailedTranslation = new(string)

	for {
		tokenType := tknzr.Next()

		switch {
		case tokenType == html.StartTagToken:
			t := tknzr.Token()

			if isBoldToken(t) {
				stack.Push("bold")
				if len(*detailedTranslation) == 0 {
					*detailedTranslation += "<b>"
				} else {
					*detailedTranslation += ", <b>"
				}
			}

		case tokenType == html.EndTagToken:
			// t := tknzr.Token()

			head, err := stack.Head()
			if err == nil && head == "bold" {
				*detailedTranslation += "</b>"
				stack.Pop()
			}
		case tokenType == html.TextToken:
			t := tknzr.Token()

			_, err := stack.Head()

			if err == nil {
				*detailedTranslation += t.Data
			}
		case tokenType == html.ErrorToken:
			return detailedTranslation, err
		}
	}
}

func isBoldToken(token html.Token) bool {
	if token.Data == "font" {
		idx := slices.IndexFunc(token.Attr, func(attr html.Attribute) bool { return attr.Key == "color" && attr.Val == "831b03" })

		return idx >= 0
	}

	return false
}

func getScarnikSuggestions(searchTerm string) ([]Suggestion, error) {
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
