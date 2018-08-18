package main

import (
	"fmt"
	"io"
	"os"
	"plugin"

	"github.com/jamesbcook/chatbot/kbchat"
)

//ActivePlugin interface for active plugins
type ActivePlugin interface {
	Get(message string) (string, error)
	Send(subscription kbchat.SubscriptionMessage, message string) error
	CMD() string
	Help() string
	Debug(bool, *io.Writer)
}

func activePluginEnvironment() string {
	return os.Getenv("CHATBOT_ACTIVE_PLUGINS")
}

func loadActivePlugins(files []string) error {
	help = make([]string, len(files))
	for x, f := range files {
		p, err := plugin.Open(f)
		if err != nil {
			return fmt.Errorf("Can't open plugin file %s %v", f, err)
		}

		apSym, err := p.Lookup("AP")
		ap := apSym.(ActivePlugin)
		help[x] = ap.Help()
		ap.Debug(debug, &writers)

		if _, ok := activePluginMap[ap.CMD()]; !ok {
			activePluginMap[ap.CMD()] = ap
		}
	}
	return nil
}
