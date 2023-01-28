package main

import (
	"github.com/gofiber/fiber/v2"
	"go-bot/handler"
	"go-bot/route"
	"go-bot/storage"
	"go-bot/task_manager"
	"log"
	"os"
)

// @title Go xaos.mobi bot
// @version 0.0.0
// @description API for xaos.mobi bot
// @host localhost:5504
// @BasePath /
func main() {
	app := fiber.New()
	conn := os.Getenv("XAOSBOT_CONNECTION_STRING")
	if conn == "" {
		log.Fatal("FATAL: Set XAOSBOT_CONNECTION_STRING to connect to the database")
	}
	accountStorage, err := storage.NewAccountStorage(conn)
	if err != nil {
		log.Fatal(err.Error())
	}
	statusStorage, err := storage.NewStatusStorage(conn)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer accountStorage.Close()
	taskManager := task_manager.New(accountStorage, statusStorage)
	controller := handler.New(taskManager)

	route.InitRoutes(app, controller)

	port := os.Getenv("XAOSBOT_PORT")
	if port == "" {
		log.Println("WARN: XAOSBOT_PORT env var is not set, using 5504 as default")
		port = "5504"
	}

	err = taskManager.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		log.Fatal(taskManager.Start())
	}()

	log.Fatal(app.Listen(":" + port))
}
