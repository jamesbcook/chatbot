package team

import (
	"testing"

	"github.com/jamesbcook/chat-bot/kbchat"
)

func TestMembers(t *testing.T) {
	api, err := kbchat.Start("team")
	if err != nil {
		t.Fatalf("Couldn't connect to team api %v", err)
	}
	res, err := Get(api, "optiv", Members)
	if err != nil {
		t.Fatalf("Error getting memberes")
	}
	if len(res) <= 0 {
		t.Fatalf("Length of members is 0 or less")
	}
}
