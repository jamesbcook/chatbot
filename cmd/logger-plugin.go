package main

//Logger is used for the Write function in the plugins
type Logger interface {
	Write(p []byte) (int, error)
}
