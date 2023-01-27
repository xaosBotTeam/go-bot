package handler

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-multierror"
	models "github.com/xaosBotTeam/go-shared-models/task"
	"go-bot/task_manager"
	"net/http"
	"strconv"
)

func New(service *task_manager.TaskManager) *BotController {
	if service == nil {
		return nil
	}

	return &BotController{service: service}
}

type BotController struct {
	service *task_manager.TaskManager
}

func (b *BotController) PutAccountTaskConfig(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return multierror.Append(err, c.SendStatus(http.StatusBadRequest))
	}
	var status models.Status
	err = json.Unmarshal(c.Body(), &status)
	if err != nil {
		return multierror.Append(err, c.SendStatus(http.StatusBadRequest))
	}

	err = b.service.UpdateStatus(id, status)
	if err != nil {
		return multierror.Append(err, c.SendStatus(http.StatusBadRequest))
	}

	return c.SendStatus(http.StatusOK)
}
