package kbchat

//ValidUser checks if the user is allowed to talk to the bot
func (api *API) ValidUser(name string) bool {
	for _, user := range api.ValidUsers {
		if user == name {
			return true
		}
	}
	return false
}
