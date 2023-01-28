package handler

import (
	"encoding/json"
	"go-bot/task_manager"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	models "github.com/xaosBotTeam/go-shared-models/task"
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

//	@Summary		Status by id
//	@Description	get status by ID
//	@ID				get-status-by-id
//
// @Accept json
// @Param id query int true "account id"
// @Param config body string true "new config"
// @Router			/task/id [put]
func (b *BotController) PutAccountTaskConfig(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	var status models.Status
	err = json.Unmarshal(c.Body(), &status)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	err = b.service.UpdateStatus(id, status)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	return c.SendStatus(http.StatusOK)
}

//	@Summary		Status by id
//	@Description	get status by ID
//	@ID				get-status-by-id
//
// @Produce			json
// @Param			id	query int false "account id"
// @Router			/task/id [get]
func (b *BotController) GetAccountStatusById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	status, err := b.service.GetStatusById(id)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	_, err = json.Marshal(status)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.JSON(status)
}

//	@Summary		All statuses
//	@Description	get all statuses
//	@ID				get-all-status
//
// @Produce			json
// @Router			/task/ [get]
func (b *BotController) GetAllStatuses(c *fiber.Ctx) error {
	statuses, err := b.service.GetAllStatuses()
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}
	_, err = json.Marshal(statuses)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(statuses)
}

//	@Summary		Restart task manager
//	@Description	restart task manager in order to add new accounts
//	@ID				restart-task-manager
//
// @Router			/refresh [get]
func (b *BotController) RestartTaskManager(c *fiber.Ctx) error {
	b.service.RefreshAccounts()
	return c.SendStatus(http.StatusOK)
}
