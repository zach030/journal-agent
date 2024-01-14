package main

import (
	"testing"
)

func TestAgent(t *testing.T) {
	agent := NewJournalAgent("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", "", "/Users/zach/Documents/flashnotes/iCloud/flashnotes")
	review, err := agent.NotesReview()
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(review)
}
