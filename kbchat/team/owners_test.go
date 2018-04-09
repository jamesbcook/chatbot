package team

import (
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

func TestOwners(t *testing.T) {
	api, err := kbchat.Start("team")
	if err != nil {
		t.Fatalf("Couldn't connect to team api %v", err)
	}
	res, err := Get(api, "optiv", Owners)
	if err != nil {
		t.Fatalf("Error getting owners")
	}
	if len(res) <= 0 {
		t.Fatalf("Length of owners is 0 or less")
	}
}
