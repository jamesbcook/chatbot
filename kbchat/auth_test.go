package kbchat_test

import (
	"testing"

	"github.com/jamesbcook/chat-bot/kbchat"
)

func TestValidUser(t *testing.T) {
	api := &kbchat.API{}
	api.ValidUsers = []string{"bob", "alice"}
	for _, user := range api.ValidUsers {
		valid := api.ValidUser(user)
		if !valid {
			t.Fatalf("Expected valid user got %v", user)
		}
	}

	if api.ValidUser("Not a valid user") {
		t.Fatalf("Expected invalid user")
	}
}
