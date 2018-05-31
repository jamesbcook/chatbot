package kbchat

import (
	"bufio"
	"io"
	"os"
)

//Senders

type sendMessageBody struct {
	Body string `json:"body,omitempty"`
}

type sendMessageOptions struct {
	ConversationID string          `json:"conversation_id,omitempty"`
	Channel        Channel         `json:"channel,omitempty"`
	Message        sendMessageBody `json:"message,omitempty"`
	Filename       string          `json:"filename,omitempty"`
	Title          string          `json:"title,omitempty"`
}

type sendMessageParams struct {
	Options sendMessageOptions `json:"options"`
}

type sendMessageArg struct {
	Method string            `json:"method"`
	Params sendMessageParams `json:"params"`
}

//Sender Information
type Sender struct {
	UID        string `json:"uid"`
	Username   string `json:"username"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
}

//Channel information
type Channel struct {
	Name        string `json:"name"`
	Public      bool   `json:"public,omitempty"`
	TopicType   string `json:"topic_type,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
	MembersType string `json:"members_type,omitempty"`
}

//Conversation information
type Conversation struct {
	ID      string  `json:"id"`
	Unread  bool    `json:"unread"`
	Channel Channel `json:"channel"`
}

//Result of conversations
type Result struct {
	Conversations []Conversation `json:"conversations"`
}

//Inbox top layer of result
type Inbox struct {
	Result Result `json:"result"`
}

//Text of a message
type Text struct {
	Body string `json:"body"`
}

//Content of a message
type Content struct {
	Type string `json:"type"`
	Text Text   `json:"text"`
}

//Message container
type Message struct {
	Content Content `json:"content"`
	Sender  Sender  `json:"sender"`
}

//MessageHolder for the keybase api
type MessageHolder struct {
	Msg Message `json:"msg"`
}

//ThreadResult from messages
type ThreadResult struct {
	Messages []MessageHolder `json:"messages"`
}

//Thread top layer of result keybaase api
type Thread struct {
	Result ThreadResult `json:"result"`
}

// API is the main object used for communicating with the Keybase JSON API
type API struct {
	Input      io.Writer
	Output     *bufio.Scanner
	ValidUsers []string
	Proc       *os.Process
	username   string
}

// SubscriptionMessage contains a message and conversation object
type SubscriptionMessage struct {
	Message      Message
	Conversation Conversation
}

// NewMessageSubscription has methods to control the background message fetcher loop
type NewMessageSubscription struct {
	newMsgsCh  <-chan SubscriptionMessage
	errorCh    <-chan error
	shutdownCh chan struct{}
}
