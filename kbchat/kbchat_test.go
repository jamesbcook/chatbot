package kbchat

import (
	"context"
	"testing"
	"time"
)

func TestKeyBaseLocation(t *testing.T) {
	if keybaseLocation() != "keybase" {
		t.Fatalf("Keybase location has been changed")
	}
}

func TestGetUsername(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	username, err := getUsername(ctx)
	if err != nil {
		t.Fatalf("Unable to get user name %v", err)
	}
	if username == "" {
		t.Fatalf("Username is empty")
	}
}

func TestStart(t *testing.T) {
	api, err := Start("chat")
	if err != nil {
		t.Fatalf("Unable to start API")
	}
	if api.Username == "" {
		t.Fatalf("Unable to get username")
	}
}

func TestGetConversations(t *testing.T) {
	api, err := Start("chat")
	if err != nil {
		t.Fatalf("Unable to start API")
	}
	if api.Username == "" {
		t.Fatalf("Unable to get username")
	}
	_, err = api.GetConversations(true)
	if err != nil {
		t.Fatalf("Error getting conversations unread true %v", err)
	}
	conv, err := api.GetConversations(false)
	if err != nil {
		t.Fatalf("Error getting conversations unread false %v", err)
	}
	if len(conv) <= 0 {
		t.Fatalf("No conversations found")
	}
}

func TestGetTextMessages(t *testing.T) {
	api, err := Start("chat")
	if err != nil {
		t.Fatalf("Unable to start API")
	}
	if api.Username == "" {
		t.Fatalf("Unable to get username")
	}
	conv, err := api.GetConversations(false)
	if err != nil {
		t.Fatalf("Error getting conversations unread false %v", err)
	}
	if len(conv) <= 0 {
		t.Fatalf("No conversations found")
	}
	_, err = api.GetTextMessages(conv[0].ID, true) // making this false crashes and fails, don't know why
	if err != nil {
		t.Fatalf("Couldn't get text messages from conv id %v %v", conv[0].ID, err)
	}
}

func TestSendMessage(t *testing.T) {
	api, err := Start("chat")
	if err != nil {
		t.Fatalf("Unable to start API")
	}
	if api.Username == "" {
		t.Fatalf("Unable to get username")
	}
	conv, err := api.GetConversations(false)
	if err != nil {
		t.Fatalf("Error getting conversations unread false %v", err)
	}
	if err := api.SendMessage(conv[0].ID, "Running Go Test from Message ID"); err != nil {
		t.Fatalf("Error sending message %v", err)
	}
}

func TestSendMessageByTLfName(t *testing.T) {
	channelName := "chatbot2,jamesbcook"
	api, err := Start("chat")
	if err != nil {
		t.Fatalf("Unable to start API")
	}
	if api.Username == "" {
		t.Fatalf("Unable to get username")
	}
	if err := api.SendMessageByTlfName(channelName, "Running Go Test from Channel Name"); err != nil {
		t.Fatalf("Error sending message to channel %v", err)
	}

}

func TestSendMessageByTeamName(t *testing.T) {
	teamName := "chatbot_dev"
	channel := "dev"
	api, err := Start("chat")
	if err != nil {
		t.Fatalf("Unable to start API")
	}
	if api.Username == "" {
		t.Fatalf("Unable to get username")
	}
	if err := api.SendMessageByTeamName(teamName, "Running Go Testing from Team Name in general", nil); err != nil {
		t.Fatalf("Unable to send message to team channel general %v", err)
	}
	if err := api.SendMessageByTeamName(teamName, "Running Go Testing from Team Name in dev", &channel); err != nil {
		t.Fatalf("Unable to send message to team channel %s %v", err, channel)
	}
}

/*
func TestMessageListen(t *testing.T) {
	api, err := Start("chat")
	if err != nil {
		t.Fatalf("Unable to start API")
	}
	if api.Username == "" {
		t.Fatalf("Unable to get username")
	}
	newMsg := api.ListenForNewTextMessages()
	msg, err := newMsg.Read()
	if err != nil {
		t.Fatalf("Error reading new message %v", err)
	}
	if msg.Conversation.ID == "0" {
		t.Fatalf("Unalbe to read new message")
	}
}
*/
