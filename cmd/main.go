package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
	"github.com/jamesbcook/print"
)

var (
	activePluginMap     = make(map[string]ActivePlugin)
	backgroundPluginMap = make(map[string]*backgroundPluginHolder)
	extraPlugins        = make(map[string]Extra)
	writers             io.Writer
	help                []string
	gitCommit           string
	binVersion          string
	debug               bool
)

type loader func([]string) error

func printFiles(files []string) error {
	if debug {
		for _, file := range files {
			print.Statusf("Loading %s\n", file)
		}
	}
	return nil
}

func decorateLoader(files []string, loaders ...loader) {
	for _, load := range loaders {
		if err := load(files); err != nil {
			print.Badln(err)
		}
	}
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

func globFiles(path string) []string {
	files, err := filepath.Glob(path + "/*.so")
	if err != nil {
		print.Badln(err)
	}
	return files
}

func wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func main() {
	d := flag.Bool("debug", false, "Print debug statements from plugins")
	v := flag.Bool("version", false, "Current Version")
	flag.Parse()
	if *v {
		fmt.Printf("chatbot v%s %s\n", binVersion, gitCommit)
		os.Exit(0)
	}
	debug = *d
	type loaderEnv struct {
		envFunc           func() string
		envWarningMessage string
		loadFunc          loader
	}

	loaders := []loaderEnv{}

	loaders = append(loaders, loaderEnv{activePluginEnvironment, "CHATBOT_ACTIVE_PLUGINS", loadActivePlugins})
	loaders = append(loaders, loaderEnv{backgroundPluginEnvironment, "CHATBOT_BACKGROUND_PLUGINS", loadBackgroundPlugins})
	loaders = append(loaders, loaderEnv{extraPluginEnvironment, "CHATBOT_EXTRA_PLUGINS", loadExtraPlugins})

	for _, load := range loaders {
		go func(l loaderEnv) {
			for {
				path := l.envFunc()
				if path == "" {
					print.Warningf("Missing %s\n", l.envWarningMessage)
					wait(30)
					continue
				}
				plugins := globFiles(path)

				decorateLoader(plugins, printFiles, l.loadFunc)
				wait(30)
			}
		}(load)
	}

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
	errorWriter := print.Error(&writers)

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
		go func(messageID, message string) {
			for _, ep := range extraPlugins {
				res, err := ep.Get(message)
				if res == "" || err != nil {
					errorWriter(err)
					return
				}
				if err := ep.Send(messageID, res); err != nil {
					errorWriter(err)
					return
				}
			}
		}(subscription.Conversation.ID, subscription.Message.Content.Text.Body)

		command := strings.Fields(subscription.Message.Content.Text.Body)
		if _, ok := activePluginMap[command[0]]; !ok {
			continue
		}
		arg := strings.Join(command[1:], " ")
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
