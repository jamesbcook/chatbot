package team

import (
	"testing"

	"github.com/jamesbcook/chat-bot/kbchat"
)

func TestAdmins(t *testing.T) {
	api, err := kbchat.Start("team")
	if err != nil {
		t.Fatalf("Couldn't connect to team api %v", err)
	}
	res, err := Get(api, "optiv", Admins)
	if err != nil {
		t.Fatalf("Error getting admins")
	}
	if len(res) <= 0 {
		t.Fatalf("Length of admins is 0 or less")
	}
}
