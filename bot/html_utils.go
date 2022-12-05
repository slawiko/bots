package main

import "golang.org/x/net/html"

func searchAttributes(attrs []html.Attribute, key string, val string) bool {
	for i := 0; i < len(attrs); i++ {
		if attrs[i].Key == key && attrs[i].Val == val {
			return true
		}
	}

	return false
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
