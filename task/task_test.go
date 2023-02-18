package task_test

import (
	"go-bot/task"
	"os"
	"testing"

	"github.com/xaosBotTeam/go-shared-models/account"
	"github.com/xaosBotTeam/go-shared-models/status"
)


type Abstract interface {
	Do(account.Account, status.Status) error
}

func testAbstractTask(testTask Abstract) error {
	testAcc := account.Account{
		URL: os.Getenv("XAOSBOT_TEST_ACC_URL"),
		Owner: 0,
	}
	return testTask.Do(testAcc, status.Status{})
}

func TestTravelling_Do(t *testing.T) {
	err := testAbstractTask(task.NewTravelling())
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestChests_Do(t *testing.T) {
	err := testAbstractTask(task.NewChests())
	if err != nil {
		t.Fatal(err.Error())
	}
}