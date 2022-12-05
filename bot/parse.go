package main

import (
	"errors"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func DetailedTranslationParse(body io.Reader) (string, error) {
	tknzr := html.NewTokenizer(body)
	stack := stack{
		stack: make([]html.Token, 0),
	}
	var builder strings.Builder
	tooLong := false

	for {
		tokenType := tknzr.Next()
		if tokenType == html.ErrorToken {
			translation := builder.String()
			translation = strings.TrimSpace(translation)
			if len(translation) == 0 {
				return "", errors.New("nothing is parsed")
			}
			if tooLong {
				translation += "\n\n<b><i>... далей чытайце на skarnik.by</i></b>"
			}
			return translation, nil
		}

		// 300 - empirical number in favor of simplicity
		if builder.Len()+300 > TelegramMessageMaxSize {
			tooLong = true
		}

		t := tknzr.Token()

		// translation in skarnik is inside p#trn element, so no need to check any other elements
		if stack.Empty() && !isPTRN(t) {
			continue
		}

		switch {
		case tokenType == html.StartTagToken:
			if tooLong {
				continue
			}

			if isItalic(t) || isGreyText(t) {
				stack.Push(t)
				builder.WriteString("<i>")
			} else if isTranslationToken(t) {
				stack.Push(t)
				builder.WriteString("<b>")
			} else if isP(t) {
				stack.Push(t)
			}
		case tokenType == html.EndTagToken:
			head, err := stack.Head()
			if err != nil {
				continue
			}

			if isBr(t) {
				builder.WriteString("\n")
			} else if isItalic(head) || isGreyText(head) {
				builder.WriteString("</i>")
				stack.Pop()
			} else if isTranslationToken(head) {
				builder.WriteString("</b>")
				stack.Pop()
			} else if isP(t) {
				stack.Pop()
			}
		case tokenType == html.TextToken:
			if stack.Empty() || tooLong {
				continue
			}

			builder.WriteString(t.Data)
		}
	}
}

func ShortTranslationParse(body io.Reader) (string, error) {
	tknzr := html.NewTokenizer(body)
	stack := stack{
		stack: make([]html.Token, 0),
	}
	var builder strings.Builder

	for {
		tokenType := tknzr.Next()
		if tokenType == html.ErrorToken {
			if builder.Len() == 0 {
				return "", errors.New("nothing is parsed")
			}
			return builder.String(), nil
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

				if builder.Len() == 0 {
					builder.WriteString("<b>")
				} else {
					builder.WriteString(", <b>")
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
				builder.WriteString("</b>")
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
				builder.WriteString(t.Data)
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
