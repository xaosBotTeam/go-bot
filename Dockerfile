FROM golang:1.19-alpine as build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY * ./

RUN go build -o xaosbot

FROM scratch

COPY --from=build /app/xaosbot /bin/xaosbot
ENTRYPOINT [ "/bin/xaosbot" ]
