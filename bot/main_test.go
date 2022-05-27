package main

import "testing"

func TestPrepareRequestText(t *testing.T) {
	t.Run("Should handle all cases", func(t *testing.T) {
		res := PrepareRequestText("Як будзе іў’'")
		if res != "як будзе ищъъ" {
			t.Errorf("expected (%s), got (%s)", "як будзе ищъъ", res)
		}
	})
}
