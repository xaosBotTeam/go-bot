package storage

import (
	"errors"
	"github.com/xaosBotTeam/go-shared-models/account"
)

func NewFakeAccountStorage() *FakeAccountStorage {
	return &FakeAccountStorage{accounts: make([]account.Account, 0)}
}

var _ (AbstractAccountStorage) = (*FakeAccountStorage)(nil)

type FakeAccountStorage struct {
	accounts []account.Account
}

func (f *FakeAccountStorage) GetAll() ([]account.Account, error) {
	copyAccounts := make([]account.Account, len(f.accounts))
	copy(copyAccounts, f.accounts)
	return copyAccounts, nil
}

func (f *FakeAccountStorage) GetById(id int) (account.Account, error) {
	for _, acc := range f.accounts {
		if acc.ID == id {
			return acc, nil
		}
	}
	return account.Account{}, errors.New("not found")
}

func (f *FakeAccountStorage) GetTable() string {
	return "FakeTable"
}

func (f *FakeAccountStorage) Close() error {
	return nil
}

func (f *FakeAccountStorage) Add(url string, ownerId int) (account.Account, error) {
	newAcc := account.Account{
		ID:           len(f.accounts),
		GameID:       0,
		FriendlyName: "New account",
		Owner:        ownerId,
		URL:          url,
		EnergyLimit:  1000,
	}

	f.accounts = append(f.accounts, newAcc)
	return newAcc, nil
}

func (f *FakeAccountStorage) Update(acc account.Account) error {
	f.accounts[acc.ID] = acc
	return nil
}
