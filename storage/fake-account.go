package storage

import (
	"errors"
	account "github.com/xaosBotTeam/go-shared-models/dbAccountInformation"
)

func NewFakeAccountStorage() *FakeAccountStorage {
	return &FakeAccountStorage{accounts: make([]account.DbAccountInformation, 0)}
}

type FakeAccountStorage struct {
	accounts []account.DbAccountInformation
}

func (f *FakeAccountStorage) GetAll() ([]account.DbAccountInformation, error) {
	copyAccounts := make([]account.DbAccountInformation, len(f.accounts))
	copy(copyAccounts, f.accounts)
	return copyAccounts, nil
}

func (f *FakeAccountStorage) GetById(id int) (account.DbAccountInformation, error) {
	for _, acc := range f.accounts {
		if acc.ID == id {
			return acc, nil
		}
	}
	return account.DbAccountInformation{}, errors.New("not found")
}

func (f *FakeAccountStorage) GetTable() string {
	return "FakeTable"
}

func (f *FakeAccountStorage) Close() error {
	return nil
}
func (f *FakeAccountStorage) Add(acc account.DbAccountInformation) {
	f.accounts = append(f.accounts, acc)
}
