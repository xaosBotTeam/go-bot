package info_collector

import (
	"github.com/xaosBotTeam/go-shared-models/account"
	"os"
	"testing"
)

func TestNicknameCollector(t *testing.T) {
	collector := Nickname{}
	acc := account.Account{
		ID:           0,
		GameID:       0,
		FriendlyName: "",
		Owner:        0,
		URL:          os.Getenv("XAOSBOT_TEST_ACC_URL"),
		EnergyLimit:  0,
	}
	acc, err := collector.Collect(acc)
	if err != nil {
		t.Fatal(err.Error())
	}
	if acc.FriendlyName != "UroborosQ" {
		t.Fail()
	}
}

func TestEnergyLimit(t *testing.T) {
	collector := EnergyLimit{}
	acc := account.Account{
		ID:           0,
		GameID:       0,
		FriendlyName: "",
		Owner:        0,
		URL:          os.Getenv("XAOSBOT_TEST_ACC_URL"),
		EnergyLimit:  0,
	}
	acc, err := collector.Collect(acc)
	if err != nil {
		t.Fatal(err.Error())
	}
	if acc.EnergyLimit <= 0 {
		t.Fail()
	}
}

func TestGameId(t *testing.T) {
	collector := GameId{}
	acc := account.Account{
		ID:           0,
		GameID:       0,
		FriendlyName: "",
		Owner:        0,
		URL:          os.Getenv("XAOSBOT_TEST_ACC_URL"),
		EnergyLimit:  0,
	}
	acc, err := collector.Collect(acc)
	if err != nil {
		t.Fatal(err.Error())
	}
	if acc.GameID != 1226207 {
		t.Fail()
	}
}
