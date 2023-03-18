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
    apk add --update bash curl udev ttf-freefont chromium chromium-chromedriver
ENV ROOT=/go/src
ENV PATH=$PATH:$GOPATH/bin

WORKDIR ${ROOT}

RUN addgroup -S dockergroup && adduser -S docker -G dockergroup
USER docker

COPY --from=build ${ROOT}/client ${ROOT}
COPY src/config/config.ini ${ROOT}
COPY src/config/moegi-settings.txt ${ROOT}
COPY src/config/self-intro.txt ${ROOT}

CMD ["./client"]
