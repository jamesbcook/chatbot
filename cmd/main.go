package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"time"

	"github.com/jamesbcook/chatbot/kbchat"
)

//Getter is used for the Get function in the plugins
type Getter interface {
	Get(string) (string, error)
}

//Sender is used for the Send function in the plugins
type Sender interface {
	Send(string, string) error
}

//Logger is used for the Write function in the plugins
type Logger interface {
	Write(p []byte) (int, error)
}

type pluginHolder struct {
	Getter
	Sender
}

type backgroundPluginHolder struct {
	plug *plugin.Plugin
	Logger
	Start    func(io.Writer)
	Validate func(string) bool
}

var (
	pluginMap           = make(map[string]*pluginHolder)
	backgroundPluginMap = make(map[string]*backgroundPluginHolder)
	writers             io.Writer
	help                []string
)

func loadActivePlugins(files []string) error {
	help = make([]string, len(files))
	for x, f := range files {
		p, err := plugin.Open(f)
		if err != nil {
			return fmt.Errorf("Can't open plugin file %s %v", f, err)
		}

		nameSym, err := p.Lookup("CMD")
		if err != nil {
			return fmt.Errorf("Can't find CMD symbol %v in %s", err, f)
		}

		h, err := p.Lookup("Help")
		if err != nil {
			return fmt.Errorf("Can't find Help symbol %v in %s", err, f)
		}

		help[x] = *h.(*string)
		plugHolder := &pluginHolder{}
		tmp, err := p.Lookup("Getter")
		if err != nil {
			return fmt.Errorf("Can't find Get symbol %v in %s", err, f)
		}

		plugHolder.Getter = tmp.(Getter)
		tmp, err = p.Lookup("Sender")
		if err != nil {
			return fmt.Errorf("Can't find Sender symbol %v in %s", err, f)
		}

		plugHolder.Sender = tmp.(Sender)
		if _, ok := pluginMap[*nameSym.(*string)]; !ok {
			pluginMap[*nameSym.(*string)] = plugHolder
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
		nameSym, err := p.Lookup("Name")
		if err != nil {
			return fmt.Errorf("Can't find Name symbol %v in %s", err, f)
		}
		if _, ok := backgroundPluginMap[*nameSym.(*string)]; !ok {
			backgroundPluginMap[*nameSym.(*string)] = &backgroundPluginHolder{plug: p}
		}
	}
	return nil
}

func errorWriter(err error) {
	output := []byte(err.Error())
	output = append(output, '\n')
	writers.Write(output)
}

func fatalErrorWriter(err error) {
	output := []byte(err.Error())
	output = append(output, '\n')
	writers.Write(output)
	os.Exit(1)
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
	go func() {
		for {
			activePluginEnv := os.Getenv("CHATBOT_ACTIVE_PLUGINS")
			if activePluginEnv == "" {
				log.Println("Missing CHATBOT_ACTIVE_PLUGINS environment variable")
			}
			activePlugins, err := filepath.Glob(activePluginEnv + "/*.so")
			if err != nil {
				log.Fatal("Error with filepath glob", err)
			}

			if err := loadActivePlugins(activePlugins); err != nil {
				log.Fatal(err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

	go func() {
		for {
			backgroundPluginEnv := os.Getenv("CHATBOT_BACKGROUND_PLUGINS")
			if backgroundPluginEnv == "" {
				log.Println("Missing CHATBOT_BACKGROUND_PLUGINS environment variable")
			}
			backgroundPlugins, err := filepath.Glob(backgroundPluginEnv + "/*.so")
			if err != nil {
				log.Fatal("Error with filepath glob", err)
			}

			if err := loadBackgroundPlugins(backgroundPlugins); err != nil {
				log.Fatal(err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

	time.Sleep(10 * time.Second)

	go func() {
		for {
			if _, ok := pluginMap["/help"]; ok {
				pluginMap["/help"].Get(strings.Join(cleanHelp(help), "\n"))
			}
			time.Sleep(30 * time.Second)
		}
	}()

	fmt.Println("Ready")

	var writerList []io.Writer
	writerList = append(writerList, os.Stdout)
	if _, ok := backgroundPluginMap["log"]; ok {
		writer, err := backgroundPluginMap["log"].plug.Lookup("Logger")
		if err != nil {
			log.Fatalf("Error looking up logger in log plugin %v", err)
		}
		backgroundPluginMap["log"].Logger = writer.(Logger)
		plugWriter := backgroundPluginMap["log"].Logger
		writerList = append(writerList, plugWriter)
	}
	writers = io.MultiWriter(writerList...)

	if _, ok := backgroundPluginMap["auth"]; ok {
		startSym, err := backgroundPluginMap["auth"].plug.Lookup("Start")
		if err != nil {
			fatalErrorWriter(fmt.Errorf("[Error] auth start symbol not found"))
		}
		//start valid user gathering in the background
		go startSym.(func(io.Writer))(writers)
		validSym, err := backgroundPluginMap["auth"].plug.Lookup("Validate")
		if err != nil {
			fatalErrorWriter(fmt.Errorf("[Error] auth validate symbol not found"))
		}
		backgroundPluginMap["auth"].Validate = validSym.(func(string) bool)
	}

	kbcRead, err := kbchat.Start("chat")
	if err != nil {
		fatalErrorWriter(fmt.Errorf("[Error] Read API: %v", err))
	}

	sub := kbcRead.ListenForNewTextMessages()
	for {
		msg, err := sub.Read()
		if err != nil {
			errorWriter(fmt.Errorf("[Error] reading message %v", err))
			continue
		}
		if _, ok := backgroundPluginMap["auth"]; ok {
			if !backgroundPluginMap["auth"].Validate(msg.Message.Sender.Username) {
				continue
			}
		}
		command := strings.Fields(msg.Message.Content.Text.Body)
		arg := strings.Join(command[1:], " ")
		errorWriter(fmt.Errorf("[Debug] %v", command))
		if _, ok := pluginMap[command[0]]; !ok {
			continue
		}
		go func() {
			res, err := pluginMap[command[0]].Get(arg)
			if err != nil {
				errorWriter(fmt.Errorf("[Error] Get command %v", err))
				res = "[Error] in command"
			}
			if len(res) <= 0 {
				return
			}
			if err := pluginMap[command[0]].Send(msg.Conversation.ID, res); err != nil {
				errorWriter(fmt.Errorf("[Error] Send command %v", err))
				return
			}
		}()
	}
}
