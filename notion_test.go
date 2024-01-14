package main

import (
	"context"
	"testing"
	"time"
)

func TestNotion_InsertNote(t *testing.T) {
	n := NewNotionAPI("secret_AKTsivZpk9O3WAx7NvEYKGXK0dSjz5l5LNVjSOmiYdB", "fd1d76d75ae1432aa21056c5dd2719a0", "a36bb2b6f1bf4641a71a5b2f0cae20c8")
	n.InsertJournal(context.Background(), []FlashNote{{"12", "23", time.Now()}})
	n.GetDatabase("a36bb2b6f1bf4641a71a5b2f0cae20c8")
}
