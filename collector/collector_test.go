package collector

import (
	"github.com/xaosBotTeam/go-shared-models/status"
	"os"
	"testing"
)

func TestNicknameCollector(t *testing.T) {
	collector := Nickname{}
	acc := status.Status{
		GameID:       0,
		FriendlyName: "",
		EnergyLimit:  0,
	}
	acc, err := collector.Collect(acc, os.Getenv("XAOSBOT_TEST_ACC_URL"))
	if err != nil {
		t.Fatal(err.Error())
	}
	if acc.FriendlyName != "UroborosQ" {
		t.Fail()
	}
}

func TestEnergyLimit(t *testing.T) {
	collector := EnergyLimit{}
	acc := status.Status{
		GameID:       0,
		FriendlyName: "",
		EnergyLimit:  0,
	}
	acc, err := collector.Collect(acc, os.Getenv("XAOSBOT_TEST_ACC_URL"))
	if err != nil {
		t.Fatal(err.Error())
	}
	if acc.EnergyLimit <= 0 {
		t.Fail()
	}
}

func TestGameId(t *testing.T) {
	collector := GameId{}
	acc := status.Status{
		GameID:       0,
		FriendlyName: "",
		EnergyLimit:  0,
	}
	acc, err := collector.Collect(acc, os.Getenv("XAOSBOT_TEST_ACC_URL"))
	if err != nil {
		t.Fatal(err.Error())
	}
	if acc.GameID != 1226207 {
		t.Fail()
	}
}
