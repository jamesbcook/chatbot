package team

import (
	"testing"

	"github.com/jamesbcook/chatbot/kbchat"
)

func TestWriters(t *testing.T) {
	api, err := kbchat.Start("team")
	if err != nil {
		t.Fatalf("Couldn't connect to team api %v", err)
	}
	res, err := Get(api, "optiv", Writers)
	if err != nil {
		t.Fatalf("Error getting writers")
	}
	if len(res) <= 0 {
		t.Fatalf("Length of writers is 0 or less")
	}
}
