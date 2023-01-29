package storage

import (
	"errors"
	account "github.com/xaosBotTeam/go-shared-models/account"
)

func NewFakeAccountStorage() *FakeAccountStorage {
	return &FakeAccountStorage{accounts: make([]account.Account, 0)}
}

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
func (f *FakeAccountStorage) Add(acc account.Account) {
	f.accounts = append(f.accounts, acc)
}
