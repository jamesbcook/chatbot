package main

import (
	"fmt"
	"io"
	"os"
	"plugin"
)

//BackgroundPlugin interface for background plugins
type BackgroundPlugin interface {
	Name() string
	Debug(bool, *io.Writer)
}

type backgroundPluginHolder struct {
	plug *plugin.Plugin
	Logger
	Authenticator
}

func backgroundPluginEnvironment() string {
	return os.Getenv("CHATBOT_BACKGROUND_PLUGINS")
}

func loadBackgroundPlugins(files []string) error {
	for _, f := range files {
		p, err := plugin.Open(f)
		if err != nil {
			return fmt.Errorf("Can't open plugin file %s %v", f, err)
		}
		bpSymbol, err := p.Lookup("BP")
		if err != nil {
			return fmt.Errorf("Can't find Name symbol %v in %s", err, f)
		}
		bp := bpSymbol.(BackgroundPlugin)
		bp.Debug(debug, &writers)
		if _, ok := backgroundPluginMap[bp.Name()]; !ok {
			backgroundPluginMap[bp.Name()] = &backgroundPluginHolder{plug: p}
		}
	}
	return nil
}
