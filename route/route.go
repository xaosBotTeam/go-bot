package route

import (
	"go-bot/handler"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

func InitRoutes(app *fiber.App, controller *handler.BotController) {
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Use(cors.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	taskApi := app.Group("/config")
	taskApi.Get("/", controller.GetAllConfigs)
	taskApi.Put("/", controller.SetConfigForAll)

	taskApi.Get("/:id", controller.GetConfigById)
	taskApi.Put("/:id", controller.UpdateConfig)

	accountApi := app.Group("/account")
	accountApi.Get("/", controller.GetAllAccounts)
	accountApi.Get("/:id", controller.GetAccountById)
	accountApi.Post("/", controller.AddAccount)
	accountApi.Delete("/:id", controller.DeleteAccount)

	statusApi := app.Group("/status")
	statusApi.Get("/", controller.GetAllStatuses)
	statusApi.Get("/:id", controller.GetStatus)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNotFound)
	})
}
