package team

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/jamesbcook/chat-bot/kbchat"
)

//Get takes a function which takes a team name
func Get(api *kbchat.API, team string, fncType func(input Results) (Output, error)) (Output, error) {
	list := fmt.Sprintf(`{"method": "list-team-memberships", "params": {"options": {"team": "%v"}}}`, team)
	if _, err := io.WriteString(api.Input, list); err != nil {
		return nil, fmt.Errorf("Writing to api input error: %v", err)
	}
	api.Output.Scan()
	var results Results
	if err := json.Unmarshal(api.Output.Bytes(), &results); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal json: %v", err)
	}
	out, err := fncType(results)
	if err != nil {
		log.Fatal(err)
	}
	return out, nil
}
