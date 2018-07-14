#!/bin/bash
#
sleep 5

# Start keybase
run_keybase

# Login
keybase oneshot

# Run chatbot
/home/keybase/chatbot/bot
