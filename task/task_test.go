package task_test

import (
	"github.com/xaosBotTeam/go-shared-models/account"
	"github.com/xaosBotTeam/go-shared-models/status"
	"go-bot/task"
	"os"
	"testing"
)

func TestTravelling(t *testing.T) {
	testAcc := account.Account{
		URL: os.Getenv("XAOSBOT_TEST_ACC_URL"),
		Owner: 0,
	}
	testTask := task.NewTravelling()
	err := testTask.Do(testAcc, status.Status{})
	if err != nil {
		t.Fatal(err.Error())
	}
}