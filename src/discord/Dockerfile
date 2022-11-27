# build用のコンテナ
FROM golang:1.18-alpine AS build
RUN apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    apk del tzdata

ENV ROOT=/go/src
WORKDIR ${ROOT}

COPY ./src ${ROOT}

RUN go mod download \
	&& CGO_ENABLED=0 GOOS=linux go build -o client ./discord

# server用のコンテナ
FROM alpine:3.15.4
RUN apk add tzdata wget curl unzip make && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    apk del tzdata && \
    apk add --update bash curl
RUN apk add --update \
    wget \
    udev \
    ttf-freefont \
    chromium \
    chromium-chromedriver
# FROM ubuntu:latest
# ENV TZ=Asia/Tokyo
# RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
# RUN apt-get update && \
#     apt-get install -y wget curl unzip make tzdata software-properties-common libnss3-dev libgdk-pixbuf2.0-dev libgtk-3-dev libxss-dev

# RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add && \
#     echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' | tee /etc/apt/sources.list.d/google-chrome.list &&\
#     apt-get update -qq &&\
#     apt-get install -y google-chrome-stable libgconf-2-4 &&\
#     CHROMEDRIVER_VERSION=`curl -sS chromedriver.storage.googleapis.com/LATEST_RELEASE` \
#     && curl -sS -o /tmp/chromedriver_linux64.zip http://chromedriver.storage.googleapis.com/$CHROMEDRIVER_VERSION/chromedriver_linux64.zip \
#     && unzip /tmp/chromedriver_linux64.zip \
#     && mv chromedriver /usr/local/bin/

ENV ROOT=/go/src
ENV PATH=$PATH:$GOPATH/bin

WORKDIR ${ROOT}

RUN addgroup -S dockergroup && adduser -S docker -G dockergroup
USER docker

COPY --from=build ${ROOT}/client ${ROOT}
COPY src/config.ini ${ROOT}
COPY src/discord/self-intro.txt ${ROOT}

CMD ["./client"]
