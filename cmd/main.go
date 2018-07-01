package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

//ActivePlugin interface for active plugins
type ActivePlugin interface {
	Get(message string) (string, error)
	Send(subscription kbchat.SubscriptionMessage, message string) error
	CMD() string
	Help() string
	Debug(bool, *io.Writer)
}

//BackgroundPlugin interface for background plugins
type BackgroundPlugin interface {
	Name() string
	Debug(bool, *io.Writer)
}

//Authenticator interface for auth plugins
type Authenticator interface {
	Start()
	Validate(user string) bool
}

//Logger is used for the Write function in the plugins
type Logger interface {
	Write(p []byte) (int, error)
}

type backgroundPluginHolder struct {
	plug *plugin.Plugin
	Logger
	Authenticator
}

var (
	activePluginMap     = make(map[string]ActivePlugin)
	backgroundPluginMap = make(map[string]*backgroundPluginHolder)
	writers             io.Writer
	errorWriter         func(v error)
	help                []string
	debug               bool
)

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

func cleanHelp(input []string) []string {
	output := []string{}
	for _, value := range input {
		if value != "" {
			output = append(output, value)
		}
	}
	return output
}

func main() {
	d := flag.Bool("debug", false, "Print debug statements from plugins")
	flag.Parse()
	debug = *d
	go func() {
		for {
			activePluginEnv := os.Getenv("CHATBOT_ACTIVE_PLUGINS")
			if activePluginEnv == "" {
				print.Warningln("Missing CHATBOT_ACTIVE_PLUGINS environment variable")
			}
			activePlugins, err := filepath.Glob(activePluginEnv + "/*.so")
			if err != nil {
				print.Badln("Error with filepath glob", err)
			}

			if err := loadActivePlugins(activePlugins); err != nil {
				print.Badln(err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

	go func() {
		for {
			backgroundPluginEnv := os.Getenv("CHATBOT_BACKGROUND_PLUGINS")
			if backgroundPluginEnv == "" {
				print.Warningln("Missing CHATBOT_BACKGROUND_PLUGINS environment variable")
			}
			backgroundPlugins, err := filepath.Glob(backgroundPluginEnv + "/*.so")
			if err != nil {
				print.Badln("Error with filepath glob", err)
			}

			if err := loadBackgroundPlugins(backgroundPlugins); err != nil {
				print.Badln(err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

	time.Sleep(10 * time.Second)

	go func() {
		for {
			if helpPlugin, ok := activePluginMap["/help"]; ok {
				helpPlugin.Get(strings.Join(cleanHelp(help), "\n"))
			}
			time.Sleep(30 * time.Second)
		}
	}()

	print.Goodln("Ready")

	var writerList []io.Writer
	writerList = append(writerList, os.Stdout)
	if logPlugin, ok := backgroundPluginMap["log"]; ok {
		writer, err := logPlugin.plug.Lookup("Logger")
		if err != nil {
			print.Badf("Error looking up logger in log plugin %v", err)
		}
		logPlugin.Logger = writer.(Logger)
		writerList = append(writerList, logPlugin.Logger)
	}
	writers = io.MultiWriter(writerList...)
	errorWriter = print.Error(&writers)

	if authPlugin, ok := backgroundPluginMap["auth"]; ok {
		validSym, err := authPlugin.plug.Lookup("Auth")
		if err != nil {
			errorWriter(fmt.Errorf("auth symbol not found"))
			os.Exit(1)
		}
		authPlugin.Authenticator = validSym.(Authenticator)
		go authPlugin.Start()
	}

	if rateLimit, ok := backgroundPluginMap["ratelimit"]; ok {
		validSym, err := rateLimit.plug.Lookup("Auth")
		if err != nil {
			errorWriter(fmt.Errorf("auth symbol not found"))
			os.Exit(1)
		}
		rateLimit.Authenticator = validSym.(Authenticator)
		go rateLimit.Start()
	}

	kbcRead, err := kbchat.Start("chat")
	if err != nil {
		errorWriter(fmt.Errorf("Read API: %v", err))
		os.Exit(1)
	}

	newMessages := kbcRead.ListenForNewTextMessages()
	for {
		subscription, err := newMessages.Read()
		if err != nil {
			errorWriter(fmt.Errorf("reading message %v", err))
			continue
		}
		if authPlugin, ok := backgroundPluginMap["auth"]; ok {
			if !authPlugin.Validate(subscription.Message.Sender.Username) {
				continue
			}
		}
		if rateLimit, ok := backgroundPluginMap["ratelimit"]; ok {
			if !rateLimit.Validate(subscription.Message.Sender.Username) {
				continue
			}
		}
		command := strings.Fields(subscription.Message.Content.Text.Body)
		arg := strings.Join(command[1:], " ")
		if _, ok := activePluginMap[command[0]]; !ok {
			continue
		}
		go func(name, arguments string) {
			res, err := activePluginMap[name].Get(arguments)
			if err != nil {
				errorWriter(fmt.Errorf("Get command %v", err))
			}
			if len(res) <= 0 {
				res = err.Error()
			}
			if err := activePluginMap[name].Send(subscription, res); err != nil {
				errorWriter(fmt.Errorf("Send command %v", err))
				return
			}
		}(command[0], arg)
	}
}
