FROM golang:1.19-alpine as build

WORKDIR /app

RUN apk update && apk add git && git clone https://github.com/xaosBotTeam/go-bot -b development


RUN cd go-bot && go get .  && go build -o xaosbot

FROM alpine:3

COPY --from=build /app/go-bot/xaosbot /xaosbot

ENTRYPOINT [ "/xaosbot" ]
