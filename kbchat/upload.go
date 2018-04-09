package kbchat

import (
	"fmt"
	"os"
)

//Upload file to keybase and delete that file
func (api *API) Upload(convID, filepath, title string) error {
	arg := sendMessageArg{
		Method: "attach",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				ConversationID: convID,
				Filename:       filepath,
				Title:          title,
			},
		},
	}
	err := api.doSend(arg)
	if err != nil {
		return fmt.Errorf("API send error %v", err)
	}
	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("Remove file error %v", err)
	}
	return nil
}
