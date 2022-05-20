package main

import (
	"testing"
)

// func BenchmarkNetHtml(b *testing.B) {
// 	file, _ := os.Open("response.html")
// 	defer file.Close()

// 	for i := 0; i < b.N; i++ {
// 		parseDetailedTranslation(file)
// 	}
// }

// func BenchmarkGoQuery(b *testing.B) {
// 	file, _ := os.Open("response.html")
// 	defer file.Close()

// 	for i := 0; i < b.N; i++ {
// 		parseSkarnikResponse(file)
// 	}
// }

func getScarnikSuggestions(word string) ([]Suggestion, error) {
	return make([]Suggestion, 0), nil
}

func Testtranslate(t *testing.T) {
	_, err := translate("word", false)

	if err == nil {
		t.Error("translate should pass err if suggestions list is empty")
	}
}
