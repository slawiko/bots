package main

import (
	"os"
	"testing"
)

type parseTestCase struct {
	name                    string
	filename                string
	expectedTranslation     string
	expectedHTMLTranslation string
	expectedError           error
}

func TestShortTranslationParse(t *testing.T) {
	tests := []parseTestCase{{
		name:                    "Should parse correctly",
		filename:                "./test_data/деревня.html",
		expectedTranslation:     "вёска",
		expectedHTMLTranslation: "<b>вёска</b>",
		expectedError:           nil,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			file, _ := os.Open(tc.filename)
			defer file.Close()
			translation, HTMLTranslation, err := ShortTranslationParse(file)
			if err != tc.expectedError {
				t.Errorf("expected (%s), got (%s)", tc.expectedError, err)
			}
			if translation != tc.expectedTranslation {
				t.Errorf("expected (%s), got (%s)", tc.expectedTranslation, translation)
			}
			if HTMLTranslation != tc.expectedHTMLTranslation {
				t.Errorf("expected (%s), got (%s)", tc.expectedHTMLTranslation, HTMLTranslation)
			}
		})
	}
}

func TestDetailedTranslationParse(t *testing.T) {
	tests := []parseTestCase{{
		name:                "Should parse correctly",
		filename:            "./test_data/деревня.html",
		expectedTranslation: "<b>вёска</b>, <i>-кі женский род</i>",
		expectedError:       nil,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			file, _ := os.Open(tc.filename)
			defer file.Close()
			translation, err := DetailedTranslationParse(file)
			if err != tc.expectedError {
				t.Errorf("expected (%s), got (%s)", tc.expectedError, err)
			}
			if translation != tc.expectedTranslation {
				t.Errorf("expected (%s), got (%s)", tc.expectedTranslation, translation)
			}
		})
	}
}

func BenchmarkDetailedTranslationParse(b *testing.B) {
	file, _ := os.Open("./test_data/слово.html")
	defer file.Close()

	for i := 0; i < b.N; i++ {
		DetailedTranslationParse(file)
	}
}
