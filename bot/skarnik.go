package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Suggestion struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

func Translate(searchTerm string, isDetailed bool) (*string, error) {
	words := strings.Fields(searchTerm)

	suggestions, err := getSkarnikSuggestions(words[0])
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(suggestions) == 0 {
		return nil, errors.New("no translation found")
	}

	resp, err := requestSkarnik(suggestions[0])
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if isDetailed {
		return DetailedTranslationParse(resp.Body)
	} else {
		return ShortTranslationParse(resp.Body)
	}
}

func DetailedTranslationParse(body io.Reader) (translation *string, err error) {
	tknzr := html.NewTokenizer(body)
	stack := stack{
		stack: make([]html.Token, 0),
	}
	translation = new(string)

	for {
		tokenType := tknzr.Next()
		if tokenType == html.ErrorToken {
			return translation, err
		}

		t := tknzr.Token()

		// translation in skarnik is inside p#trn element, so no need to check any other elements
		if stack.Empty() && !isPTRN(t) {
			continue
		}

		switch {
		case tokenType == html.StartTagToken:
			if isItalic(t) || isGreyText(t) {
				stack.Push(t)
				*translation += "<i>"
			} else if isTranslationToken(t) {
				stack.Push(t)
				*translation += "<b>"
			} else if isP(t) {
				stack.Push(t)
			}
		case tokenType == html.EndTagToken:
			head, err := stack.Head()
			if err != nil {
				continue
			}

			if isBr(t) {
				*translation += "\n"
			} else if isItalic(head) || isGreyText(head) {
				*translation += "</i>"
				stack.Pop()
			} else if isTranslationToken(head) {
				*translation += "</b>"
				stack.Pop()
			} else if isP(t) {
				stack.Pop()
			}
		case tokenType == html.TextToken:
			if stack.Empty() {
				continue
			}

			*translation += t.Data
		}
	}
}

func ShortTranslationParse(body io.Reader) (translation *string, err error) {
	tknzr := html.NewTokenizer(body)
	stack := stack{
		stack: make([]html.Token, 0),
	}
	translation = new(string)

	for {
		tokenType := tknzr.Next()
		if tokenType == html.ErrorToken {
			return translation, err
		}

		t := tknzr.Token()

		// translation in skarnik is inside p#trn element, so no need to check any other elements
		if stack.Empty() && !isPTRN(t) {
			continue
		}

		switch {
		case tokenType == html.StartTagToken:
			if isTranslationToken(t) {
				stack.Push(t)

				if len(*translation) == 0 {
					*translation += "<b>"
				} else {
					*translation += ", <b>"
				}
			} else if isP(t) {
				stack.Push(t)
			}
		case tokenType == html.EndTagToken:
			head, err := stack.Head()
			if err != nil {
				continue
			}

			if isTranslationToken(head) {
				*translation += "</b>"
				stack.Pop()
			} else if isP(t) {
				stack.Pop()
			}
		case tokenType == html.TextToken:
			head, err := stack.Head()
			if err != nil {
				continue
			}
			if isTranslationToken(head) {
				*translation += t.Data
			}
		}
	}
}

func isTranslationToken(token html.Token) bool {
	return searchAttributes(token.Attr, "color", "831b03")
}

func isPTRN(token html.Token) bool {
	return searchAttributes(token.Attr, "id", "trn")
}

func searchAttributes(attrs []html.Attribute, key string, val string) bool {
	for i := 0; i < len(attrs); i++ {
		if attrs[i].Key == key && attrs[i].Val == val {
			return true
		}
	}

	return false
}

func isGreyText(token html.Token) bool {
	return searchAttributes(token.Attr, "color", "5f5f5f")
}

func isP(token html.Token) bool {
	return token.Data == "p"
}

func isBr(token html.Token) bool {
	return token.Data == "br"
}

func isItalic(token html.Token) bool {
	return token.Data == "i"
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
