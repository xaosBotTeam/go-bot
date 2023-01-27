FROM golang:1.19-alpine as build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o xaosbot

FROM alpine:3

COPY --from=build /app/xaosbot /xaosbot

ENTRYPOINT [ "/xaosbot" ]
