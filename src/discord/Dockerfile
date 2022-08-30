# build用のコンテナ
FROM golang:1.18-alpine AS build

ENV ROOT=/go/src/project
WORKDIR ${ROOT}

COPY ./src ${ROOT}

RUN go mod download \
	&& CGO_ENABLED=0 GOOS=linux go build -o client ./discord

# server用のコンテナ
FROM alpine:3.15.4

ENV ROOT=/go/src/project
WORKDIR ${ROOT}

RUN addgroup -S dockergroup && adduser -S docker -G dockergroup
USER docker

COPY --from=build ${ROOT}/client ${ROOT}
COPY src/config.ini ${ROOT}
COPY src/discord/self-intro.txt ${ROOT}

CMD ["./client"]
