package main

import "testing"

type TestCase struct {
	name     string
	input    string
	expected string
}

func TestPrepareRequestText(t *testing.T) {
	tests := []TestCase{
		{
			"Should handle capitalized letter",
			"Як будзе купальник",
			"як будзе купальник",
		},
		{
			"Should handle і letter",
			"олівка",
			"оливка",
		},
		{
			"Should handle ў letter",
			"ўавель",
			"щавель",
		},
		{
			"Should handle ъ letter",
			"’'",
			"ъъ",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := PrepareRequestText(test.input)
			if res != test.expected {
				t.Errorf("expected (%s), got (%s)", test.expected, res)
			}
		})
	}
}
