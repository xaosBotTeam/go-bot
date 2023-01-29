package http_bridge

import (
	"bytes"
	"fmt"
	"go-bot/random"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var mode = os.Getenv("XAOSBOT_MODE")
var path = os.Getenv("XAOSBOT_LOG_PATH")
var last = 0
var current = 0

func GetBodyBytes(url string) (io.Reader, error) {
	time.Sleep(random.RandomWaitTime())
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if mode == "DEBUG" {
		file, err := os.OpenFile(path+string(os.PathSeparator)+fmt.Sprintf("bot-page-%d.html", current), os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Printf("WARN: logging pages: %s", err.Error())
		}
		var buffer bytes.Buffer
		tee := io.TeeReader(rsp.Body, &buffer)
		_, err = io.Copy(file, tee)
		if err != nil {
			log.Printf("WARN: logging pages: %s", err.Error())
		}
		if current-last >= 10 {
			err = os.Remove(path + string(os.PathSeparator) + fmt.Sprintf("bot-page-%d.html", last))
			if err != nil {
				log.Printf("WARN: logging pages: %s", err.Error())
			}
			last++
		}
		current++
		return &buffer, nil
	}

	return rsp.Body, nil
}
