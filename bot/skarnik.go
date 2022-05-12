package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func translate(searchTerm string) (*string, *[]Suggestion, error) {
	cleanSearchTerm := cleanTerm(searchTerm)
	words := strings.Fields(cleanSearchTerm)

	if len(words) == 1 {
		return translateWord(words[0])
	}

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

		wordTranslation, err := parseSkarnikResponse(resp)
		if err != nil {
			return nil, nil, err
		}

		if index == 0 {
			translation = *wordTranslation
		} else {
			translation = translation+";\n"+*wordTranslation
		}
	}

	if len(translation) == 0 {
		return nil, nil, errors.New("No translation found")
	}

	return &translation, nil, nil
}

func translateWord(word string) (*string, *[]Suggestion, error) {
	suggestions, err := getScarnikSuggestions(word)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	if len(suggestions) == 0 {
		return nil, nil, errors.New("No translation found")
	}

	log.Println("her", word, suggestions[0])

	if word != suggestions[0].Label {
		return nil, &suggestions, nil
	}

	resp, err := requestSkarnik(suggestions[0])
	if err != nil {
		return nil, nil, err
	}

	wordTranslation, err := parseSkarnikResponse(resp)
	if err != nil {
		return nil, nil, err
	}

	return wordTranslation, nil, nil
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
