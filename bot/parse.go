package main

import (
	"golang.org/x/net/html"
	"io"
)

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

func isGreyText(token html.Token) bool {
	return searchAttributes(token.Attr, "color", "5f5f5f")
}
