package main

//Authenticator interface for auth plugins
type Authenticator interface {
	Start()
	Validate(user string) bool
}
