package main

import (
	"fmt"
	"io"
	"os"
	"plugin"
)

//Extra is a set of plugins that may act on user input, but we don't want them to be a command
type Extra interface {
	Name() string
	Get(msg string) (string, error)
	Send(messageID, msg string) error
	Debug(bool, *io.Writer)
}

func extraPluginEnvironment() string {
	return os.Getenv("CHATBOT_EXTRA_PLUGINS")
}

func loadExtraPlugins(files []string) error {
	//Need either clear extra plugins or check if it is new
	for _, f := range files {
		p, err := plugin.Open(f)
		if err != nil {
			return fmt.Errorf("Can't open plugin file %s %v", f, err)
		}
		extraSym, err := p.Lookup("Extra")
		if err != nil {
			return fmt.Errorf("Can't find Name symbol %v in %s", err, f)
		}
		ep := extraSym.(Extra)
		ep.Debug(debug, &writers)
		if _, ok := extraPlugins[ep.Name()]; !ok {
			extraPlugins[ep.Name()] = ep
		}
	}
	return nil
}
