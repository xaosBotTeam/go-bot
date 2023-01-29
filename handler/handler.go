package handler

import (
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"go-bot/task_manager"
	"log"
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

//	@Summary		Update status by id
//	@Description	get status by ID
//	@ID				get-status-by-id
//
// @Tags 			Task
// @Accept 			json
// @Param 			id path int true "account id"
// @Param 			config body models.Status true "new config"
// @Router			/task/{id} [put]
func (b *BotController) PutAccountTaskConfig(c *fiber.Ctx) error {
	c.Accepts("application/json")
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Can't parse account id")
	}
	var status models.Status
	err = json.Unmarshal(c.Body(), &status)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Can't parse body")
	}

	err = b.service.UpdateStatus(id, status)
	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).SendString("Can't find account with such id")
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Something doesn't work")
	}

	return c.SendStatus(http.StatusOK)
}

// @Summary		Get status by id
// @Description	get status by ID
// @ID				get-status-by-id
// @Tags 			Task
// @Produce		json
// @Param			id path int false "account id"
// @Router			/task/{id} [get]
func (b *BotController) GetAccountStatusById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Can't parse account id")
	}
	status, err := b.service.GetStatusById(id)
	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).SendString("Can't find account with such id")
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Something doesn't work")
	}
	return c.JSON(status)
}

// @Summary		Get all statuses
// @Description	get all statuses
// @ID				get-all-status
//
// @Tags			Task
// @Produce		json
// @Router			/task/ [get]
func (b *BotController) GetAllStatuses(c *fiber.Ctx) error {
	statuses, err := b.service.GetAllStatuses()
	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).SendString("There are no statuses in storage, please, add at least one")
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Something doesn't work")
	}
	return c.JSON(statuses)
}

// @Summary		Restart task manager
// @Description	restart task manager in order to add new accounts
// @ID				restart-task-manager
// @Tags			General
//
// @Router			/refresh [get]
func (b *BotController) RestartTaskManager(c *fiber.Ctx) error {
	b.service.RefreshAccounts()
	return c.Status(http.StatusOK).SendString("Refresh request sent to task manager")
}

// @Summary		Add new game account
// @ID				add-new-game-account
//
// @Tags 			Account
// @Produce		json
// @Param			url	query string true "account url"
// @Param			owner query int true "id of account`s owner"
// @Router			/account/ [post]
func (b *BotController) AddAccount(c *fiber.Ctx) error {
	c.Accepts("application/json")
	url := c.Query("url")
	ownerId, err := strconv.Atoi(c.Query("owner"))
	if url == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Url parameter is empty")
	} else if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Can't parse ownerId parameter")
	}

	acc, err := b.service.AddAccount(url, ownerId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Can't add account")
	}

	return c.JSON(acc)
}

// @Summary		Get game account by id
// @ID				get-account-by-id
//
// @Tags 			Account
// @Produce		json
// @Param			id	path int true "account id"
// @Router			/account/{id} [get]
func (b *BotController) GetAccountById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Can't parse account id")
	}
	acc, err := b.service.GetAccountById(id)
	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).SendString("Can't find account with such id")
	} else if err != nil {
		log.Printf(err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString("Something doesn't work")
	}
	return c.JSON(acc)
}

//	@Summary		Get all game accounts
//	@ID				get-all-accounts
//
// @Tags 			Account
// @Produce			json
// @Router			/account/ [get]
func (b *BotController) GetAllAccounts(c *fiber.Ctx) error {
	accs, err := b.service.GetAllAccounts()
	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).SendString("There are no accounts in storage, please, add at least one")
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Something doesn't work")
	}
	return c.JSON(accs)
}
