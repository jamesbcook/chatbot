# docker-chatbot

> A docker image for Chatbot

## Description

This is a quick way to deploy a [Chatbot](https://github.com/jamesbcook/chatbot) installation on your local machine.


## Usage

#### Building the image

```bash
git clone https://github.com/jamesbcook/chatbot
cd docker
nano id_chatbot (Add github ssh key)
docker build -t "yourname/chatbot:yourtag" .
```
To run as a daemon:

```bash
docker run -d --name chatbot yourname/chatbot:yourtag
```