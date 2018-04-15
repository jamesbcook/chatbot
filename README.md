# ChatBot

## Installation

Make sure to [install Keybase](https://keybase.io/download).

```bash
git clone https://github.com/jamesbcook/chat-bot.git
cd chat-bot
git submodule init
git submodule update
```

or

```bash
git clone --recurse-submodules https://github.com/jamesbcook/chat-bot.git
cd chat-bot
```

The chat-bot-plugins directory are all the plugins that can be used with chatbot.  Running the following commands will build all the plugins and the main binary:

```bash
make plugin-setup
make plugin-build
make build
```

The Following environmental variables need to be set:

* CHATBOT_ACTIVE_PLUGINS
  * ```export CHATBOT_ACTIVE_PLUGINS=/home/keybase/active-plugins/```
* CHATBOT_BACKGROUND_PLUGINS
  * ```export CHATBOT_BACKGROUND_PLUGINS=/home/keybase/background-plugins/```

## Please read the plugin README files as they may require their own environmental variables