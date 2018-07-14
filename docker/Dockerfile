# this is our first build stage, it will not persist in the final image
FROM golang:1.10.0 as builder

# Add SSH key for Private Repo (Remove when we go public)
RUN mkdir /root/.ssh/
ADD id_chatbot /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa && git config --global url."git@github.com:".insteadOf "https://github.com/" \
&& echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config

# Clone GO Repos
RUN mkdir /go/src/github.com && mkdir /go/src/github.com/jamesbcook
WORKDIR /go/src/github.com/jamesbcook
RUN git clone --recurse-submodules https://github.com/jamesbcook/chatbot.git \
	&& git clone https://github.com/jamesbcook/chatbot-plugins.git \
	&& git clone https://github.com/jamesbcook/chatbot-external-api.git 
RUN go get github.com/jamesbcook/print \
	&& go get github.com/sanzaru/go-giphy \
	&& go get golang.org/x/crypto/sha3 \
	&& go get gopkg.in/ns3777k/go-shodan.v3/shodan \
	&& go get github.com/golang/protobuf/proto \
	&& go get golang.org/x/sys/cpu

# Build Chatbot
WORKDIR $GOPATH/src/github.com/jamesbcook/chatbot
RUN make build
RUN make plugin-setup
RUN make plugin-build

# Second Stage Keybase

FROM ralphte/keybase

# Install Chrome
USER root
RUN apt-get update && apt-get install wget
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add -
RUN echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' | tee /etc/apt/sources.list.d/google-chrome.list
RUN apt-get update && apt-get install google-chrome-stable -y

USER keybase
# Set Enverment Varables
ENV CHATBOT_ACTIVE_PLUGINS=/home/keybase/chatbot/active-plugins/
ENV CHATBOT_BACKGROUND_PLUGINS=/home/keybase/chatbot/background-plugins/

RUN mkdir /home/keybase/chatbot
WORKDIR /home/keybase/chatbot
COPY --from=builder /go/src/github.com/jamesbcook/chatbot/bin .
COPY docker-entrypoint.sh /home/keybase/
USER root
RUN chmod +x /home/keybase/docker-entrypoint.sh
RUN chown keybase:keybase /home/keybase/docker-entrypoint.sh
USER keybase

ENTRYPOINT ["/home/keybase/docker-entrypoint.sh"]