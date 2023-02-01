package handler

import (
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/xaosBotTeam/go-shared-models/account"
	"github.com/xaosBotTeam/go-shared-models/config"
	"go-bot/task_manager"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func New(service *task_manager.TaskManager) *BotController {
	if service == nil {
		return nil
	}

	return &BotController{service: service}
}


func (b *BotController) newReply() map[string]string {
	reply := make(map[string]string)
	reply["data"] = ""
	reply["error"] = ""
	return reply
}

type BotController struct {
	service *task_manager.TaskManager
}

//	@Summary		Update config by id
//	@Description	get config by ID
//	@ID				get-config-by-id
//
// @Tags 			Config
// @Accept 			json
// @Param 			id path int true "account id"
// @Param 			config body config.Config true "new config"
// @Router			/config/{id} [put]
func (b *BotController) UpdateConfig(c *fiber.Ctx) error {
	c.Accepts("application/json")
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Can't parse account id")
	}
	var status config.Config
	err = json.Unmarshal(c.Body(), &status)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Can't parse body")
	}

	err = b.service.UpdateConfig(id, status)
	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON("Can't find account with such id")
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Something doesn't work")
	}

	return c.SendStatus(http.StatusOK)
}

//	@Summary		Get config by id
//	@Description	get config by ID
//	@ID				get-config-by-id
//
//	@Tags 			Config
//	@Produce		json
//	@Param			id path int false "account id"
//	@Router			/config/{id} [get]
func (b *BotController) GetConfigById(c *fiber.Ctx) error {
	reply := b.newReply()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		reply["error"] = "Can't parse account id"
		return c.Status(fiber.StatusBadRequest).JSON(reply)
	}
	configuration, err := b.service.ConfigById(id)
	if err == pgx.ErrNoRows {
		reply["error"] = "Can't find account with such id"
		return c.Status(fiber.StatusNotFound).JSON(reply)
	} else if err != nil {
		reply["error"] = "Something doesn't work"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}

	jsonStr, err := json.Marshal(configuration)
	if err != nil {
		reply["error"] = "Can't transform to JSON"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	reply["data"] = string(jsonStr)
	return c.JSON(reply)
}

//	@Summary		Get all configs
//	@Description	get all configs
//	@ID				get-all-configs
//
//	@Tags			Config
//	@Produce		json
//	@Router			/config/ [get]
func (b *BotController) GetAllConfigs(c *fiber.Ctx) error {
	reply := b.newReply()
	statuses, err := b.service.AllConfigs()
	if err == pgx.ErrNoRows {
		reply["error"] = "There are no statuses in storage, please, add at least one"
		return c.Status(fiber.StatusNotFound).JSON(reply)
	} else if err != nil {
		reply["error"] = "Something doesn't work"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	
	jsonStr, err := json.Marshal(statuses)
	reply["data"] = string(jsonStr)
	
	return c.JSON(reply)
}

//	@Summary		Restart task manager
//	@Description	restart task manager in order to add new accounts
//	@ID				restart-task-manager
//	@Tags			General
//
//	@Router			/refresh [patch]
func (b *BotController) RestartTaskManager(c *fiber.Ctx) error {
	b.service.RefreshAccounts()
	return c.Status(http.StatusOK).JSON("Refresh request sent to task manager")
}

//	@Summary		Add new game account
//	@ID				add-new-game-account
//
//	@Tags 			Account
//	@Produce		json
//	@Param			account	body account.Account true "account url"
//	@Router			/account/ [post]
func (b *BotController) AddAccount(c *fiber.Ctx) error {
	c.Accepts("application/json")
	var acc account.Account
	reply := b.newReply()
	err := c.BodyParser(&acc)
	if err != nil {
		reply["error"] = "Can't transform respond to JSON"
		return c.Status(http.StatusInternalServerError).JSON(reply)
	}
	id, err := b.service.AddAccount(acc)
	if err != nil {
		reply["error"] = "Can't add account"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	reply["data"] = strconv.Itoa(id)
	return c.JSON(acc)
}

//	@Summary		Get game account by id
//	@ID				get-account-by-id
//
//	@Tags 			Account
//	@Produce		json
//	@Param			id	path int true "account id"
//	@Router			/account/{id} [get]
func (b *BotController) GetAccountById(c *fiber.Ctx) error {
	reply := b.newReply()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		reply["error"] = "Can't parse account id"
		return c.Status(fiber.StatusBadRequest).JSON(reply)
	}
	acc, err := b.service.GetAccountById(id)
	if err == pgx.ErrNoRows {
		reply["error"] = "Can't find account with such id"
		return c.Status(fiber.StatusNotFound).JSON(reply)
	} else if err != nil {
		log.Printf(err.Error())
		reply["error"] = "Something doesn't work"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	jsonStr, err := json.Marshal(acc)
	if err != nil {
		reply["error"] = "Something doesn't work"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	reply["data"] = string(jsonStr)
	return c.JSON(reply)
}

//	@Summary		Get all game accounts
//	@ID				get-all-accounts
//
//	@Tags 			Account
//	@Produce		json
//	@Router			/account/ [get]
func (b *BotController) GetAllAccounts(c *fiber.Ctx) error {
	reply := b.newReply()
	accounts, err := b.service.GetAllAccounts()
	if err == pgx.ErrNoRows {
		reply["error"] = "There are no accounts in storage"
		return c.Status(fiber.StatusNotFound).JSON(reply)
	} else if err != nil {
		reply["error"] = "Something doesn't work"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	jsonStr, err := json.Marshal(accounts)
	if err != nil {
		reply["error"] = "Something doesn't work"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	reply["data"] = string(jsonStr)
	return c.JSON(reply)
}

//	@Summary		Update config for all
//	@Description	get config for all
//	@ID				update-config-for-all
//
// @Tags 			Config
// @Accept 			json
// @Param 			config body config.Config true "new config"
// @Router			/config/ [put]
func (b *BotController) SetConfigForAll(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var status config.Config
	err := json.Unmarshal(c.Body(), &status)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Can't parse body")
	}

	err = b.service.SetConfigForAllAccount(status)
	if err == pgx.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON("Can't find account with such id")
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON("Something doesn't work")
	}

	return c.SendStatus(http.StatusOK)
}

//	@Summary		Get status by id
//	@ID				get-status-by-id
//
// @Tags 			Status
// @Produce			json
// @Param			id path int true "account id"
// @Router			/status/{id} [get]
func (b *BotController) GetStatus(c *fiber.Ctx) error {
	reply := b.newReply()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		reply["error"] = "Can't parse account id"
		return c.Status(fiber.StatusBadRequest).JSON(reply)
	}
	stat, err := b.service.GetStatus(id)
	if err == pgx.ErrNoRows {
		reply["error"] = "Can't find any statuses"
		return c.Status(fiber.StatusNotFound).JSON(reply)
	} else if err != nil {
		reply["error"] = "Something doesn't work"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	jsonStr, err := json.Marshal(stat)
	if err != nil {
		reply["error"] = "Something doesn't work"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	reply["data"] = string(jsonStr)
	return c.JSON(reply)
}

//	@Summary		Get all statuses
//	@ID				get-statuses
//
// @Tags 			Status
// @Produce			json
// @Router			/status/ [get]
func (b *BotController) GetAllStatuses(c *fiber.Ctx) error {
	reply := b.newReply()
	stat, err := b.service.GetAllStatuses()
	if err == pgx.ErrNoRows {
		reply["error"] = "Can't find any statuses"
		return c.Status(fiber.StatusNotFound).JSON(reply)
	} else if err != nil {
		reply["error"] = "Something doesn't work"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	jsonStr, err := json.Marshal(stat)
	if err != nil {
		reply["error"] = "Something doesn't work"
		return c.Status(fiber.StatusInternalServerError).JSON(reply)
	}
	reply["data"] = string(jsonStr)
	return c.JSON(reply)
}