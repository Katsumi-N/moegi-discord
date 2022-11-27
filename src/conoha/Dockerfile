# build用のコンテナ
FROM golang:1.18-alpine AS build
RUN apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    apk del tzdata
    
ENV ROOT=/go/src/project
WORKDIR ${ROOT}
COPY ./src ${ROOT}
RUN go mod download \
	&& CGO_ENABLED=0 GOOS=linux go build -o server ./conoha

# server用のコンテナ
FROM alpine:3.15.4
RUN apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    apk del tzdata
ENV ROOT=/go/src/project
WORKDIR ${ROOT}

RUN addgroup -S dockergroup && adduser -S docker -G dockergroup
USER docker

COPY --from=build ${ROOT}/server ${ROOT}
COPY src/config.ini ${ROOT}

EXPOSE 8080
CMD ["./server"]
