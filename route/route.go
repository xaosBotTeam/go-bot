package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	_ "go-bot/docs"
	"go-bot/handler"
	"net/http"
)

func InitRoutes(app *fiber.App, controller *handler.BotController) {
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Use(cors.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Patch("/refresh", controller.RestartTaskManager)

	taskApi := app.Group("/task")
	taskApi.Get("/", controller.GetAllStatuses)
	taskApi.Get("/:id", controller.GetAccountStatusById)
	taskApi.Put("/:id", controller.PutAccountTaskConfig)

	accountApi := app.Group("/account")
	accountApi.Get("/", controller.GetAllAccounts)
	accountApi.Get("/:id", controller.GetAccountById)
	accountApi.Post("/", controller.AddAccount)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNotFound)
	})
}
