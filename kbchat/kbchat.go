package kbchat

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

//Call keybase from env
func keybaseLocation() string {
	return "keybase"
}

func getUsername(ctx context.Context) (username string, err error) {
	p := exec.Command(keybaseLocation(), "status")
	output, err := p.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("Unable to get stdout pipe %v", err)
	}
	if err = p.Start(); err != nil {
		return "", fmt.Errorf("Unable to start process %v", err)
	}

	scanner := bufio.NewScanner(output)
	if !scanner.Scan() {
		return "", errors.New("unable to find Keybase username")
	}
	toks := strings.Fields(scanner.Text())
	if len(toks) != 2 {
		return "", errors.New("invalid Keybase username output")
	}
	username = toks[1]

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("unable to run Keybase command, %v", ctx.Err())
	default:
		return username, nil
	}
}

// Start fires up the Keybase JSON API in stdin/stdout mode
func Start(chatType string) (*API, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Get username first
	username, err := getUsername(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to get username %v", err)
	}

	p := exec.Command(keybaseLocation(), chatType, "api")
	input, err := p.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to get stdin pipe %v", err)
	}
	output, err := p.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable to get stdout pipe %v", err)
	}
	if err := p.Start(); err != nil {
		return nil, fmt.Errorf("Unable to start process %v", err)
	}

	boutput := bufio.NewScanner(output)
	return &API{
		Input:    input,
		Output:   boutput,
		Username: username,
	}, nil
}

// GetConversations reads all conversations from the current user's inbox. Optionally
// can filter for unread only.
func (a *API) GetConversations(unreadOnly bool) ([]Conversation, error) {
	list := fmt.Sprintf(`{"method":"list", "params": { "options": { "unread_only": %v}}}`, unreadOnly)
	if _, err := io.WriteString(a.Input, list); err != nil {
		return nil, fmt.Errorf("io write GetConversations err %v", err)
	}
	a.Output.Scan()

	var inbox Inbox
	inboxRaw := a.Output.Bytes()
	if err := json.Unmarshal(inboxRaw, &inbox); err != nil {
		return nil, fmt.Errorf("Json unmarshal GetConversations err %v", err)
	}
	return inbox.Result.Convs, nil
}

// GetTextMessages fetches all text messages from a given conversation ID. Optionally can filter
// ont unread status.
func (a *API) GetTextMessages(convID string, unreadOnly bool) ([]Message, error) {
	read := fmt.Sprintf(`{"method": "read", "params": {"options": {"conversation_id": "%s", "unread_only": %v}}}`, convID, unreadOnly)
	if _, err := io.WriteString(a.Input, read); err != nil {
		return nil, fmt.Errorf("Unable to write to std in GetTextMessage %v", err)
	}
	a.Output.Scan()

	var thread Thread
	rawThread := a.Output.Bytes()
	if err := json.Unmarshal(rawThread, &thread); err != nil {
		return nil, fmt.Errorf("unable to decode thread: %s", err.Error())
	}

	var res []Message
	for _, msg := range thread.Result.Messages {
		if msg.Msg.Content.Type == "text" {
			res = append(res, msg.Msg)
		}
	}

	return res, nil
}

func (a *API) doSend(arg sendMessageArg) error {
	bArg, err := json.Marshal(arg)
	if err != nil {
		return fmt.Errorf("Unable to marshal json in dosend %v", err)
	}
	if _, err := io.WriteString(a.Input, string(bArg)); err != nil {
		return fmt.Errorf("Unable to write string in dosend %v", err)
	}
	if err := a.Output.Scan(); err != true {
		return fmt.Errorf("Scan error %v", err)
	}
	return nil
}

// SendMessage sends a new text message on the given conversation ID
func (a *API) SendMessage(convID string, body string) error {
	arg := sendMessageArg{
		Method: "send",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				ConversationID: convID,
				Message: sendMessageBody{
					Body: body,
				},
			},
		},
	}
	return a.doSend(arg)
}

// SendMessageByTlfName sends a message on the given TLF name
func (a *API) SendMessageByTlfName(tlfName string, body string) error {
	arg := sendMessageArg{
		Method: "send",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				Channel: Channel{
					Name: tlfName,
				},
				Message: sendMessageBody{
					Body: body,
				},
			},
		},
	}
	return a.doSend(arg)
}

//SendMessageByTeamName via the keybase api
func (a *API) SendMessageByTeamName(teamName string, body string, inChannel *string) error {
	channel := "general"
	if inChannel != nil {
		channel = *inChannel
	}
	arg := sendMessageArg{
		Method: "send",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				Channel: Channel{
					MembersType: "team",
					Name:        teamName,
					TopicName:   channel,
				},
				Message: sendMessageBody{
					Body: body,
				},
			},
		},
	}
	return a.doSend(arg)
}

// Read blocks until a new message arrives
func (m NewMessageSubscription) Read() (SubscriptionMessage, error) {
	select {
	case msg := <-m.newMsgsCh:
		return msg, nil
	case err := <-m.errorCh:
		return SubscriptionMessage{}, fmt.Errorf("Error in error read channel %v", err)
	}
}

// Shutdown terminates the background process
func (m NewMessageSubscription) Shutdown() {
	m.shutdownCh <- struct{}{}
}

//GetUnreadMessagesFromConvs via the keybase api
func (a *API) GetUnreadMessagesFromConvs(convs []Conversation) ([]SubscriptionMessage, error) {
	var res []SubscriptionMessage
	for _, conv := range convs {
		msgs, err := a.GetTextMessages(conv.ID, true)
		if err != nil {
			return nil, fmt.Errorf("GetTextMessage err %v", err)
		}
		for _, msg := range msgs {
			res = append(res, SubscriptionMessage{
				Message:      msg,
				Conversation: conv,
			})
		}
	}
	return res, nil
}

// ListenForNewTextMessages fires off a background loop to fetch new unread messages.
func (a *API) ListenForNewTextMessages() NewMessageSubscription {
	newMsgCh := make(chan SubscriptionMessage, 100)
	errorCh := make(chan error, 100)
	shutdownCh := make(chan struct{})
	sub := NewMessageSubscription{
		newMsgsCh:  newMsgCh,
		shutdownCh: shutdownCh,
		errorCh:    errorCh,
	}
	go func() {
		for {
			select {
			case <-shutdownCh:
				return
			case <-time.After(2 * time.Second):
				// Get all unread convos
				convs, err := a.GetConversations(true)
				if err != nil {
					errorCh <- err
					continue
				}
				// Get unread msgs from convs
				msgs, err := a.GetUnreadMessagesFromConvs(convs)
				if err != nil {
					errorCh <- err
					continue
				}
				// Send all the new messages out
				for _, msg := range msgs {
					newMsgCh <- msg
				}
			}
		}
	}()
	return sub
}
