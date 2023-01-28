package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	_ "go-bot/docs"
	"go-bot/handler"
	"net/http"
)

func InitRoutes(app *fiber.App, controller *handler.BotController) {
	// group for setting and getting tasks for accounts
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Get("/swagger/*", swagger.HandlerDefault) // default
	api := app.Group("/task")
	// get account task status by id
	api.Get("/id", controller.GetAccountStatusById)
	// set task config
	api.Put("/id/", controller.PutAccountTaskConfig)
	api.Get("/", controller.GetAllStatuses)

	app.Patch("/refresh", controller.RestartTaskManager)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNotFound)
	})
}
