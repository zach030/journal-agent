package main

import (
	"testing"
)

func TestAgent(t *testing.T) {
	agent := NewJournalAgent("notion-sk", "", "path")
	review, err := agent.NotesReview()
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(review)
}
