package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/xaosBotTeam/go-shared-models/account"
)

func NewAccountStorage(connString string) (AbstractAccountStorage, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	table := "bot.Accounts"
	createSchemaString := `CREATE SCHEMA IF NOT EXISTS bot`
	createTableString := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id  		  int PRIMARY KEY NOT NULL,
	game_id       int NOT NULL,
	friendly_name text NOT NULL,
	owner_id      int NOT NULL,
	url           text NOT NULL,
	energy_limit  int NOT NULL)`, table)

	_, err = conn.Exec(context.Background(), createSchemaString)
	if err != nil {
		return nil, err
	}
	_, err = conn.Exec(context.Background(), createTableString)
	if err != nil {
		return nil, err
	}

	return &AccountStorage{db: conn,
		table: table}, nil
}

type AbstractAccountStorage interface {
	GetAll() ([]account.Account, error)
	GetById(id int) (account.Account, error)
	GetTable() string
	Close() error
}

type AccountStorage struct {
	db    *pgx.Conn
	table string
}

func (a *AccountStorage) GetAll() ([]account.Account, error) {
	row := a.db.QueryRow(context.Background(), fmt.Sprintf(`SELECT COUNT(*) FROM %s`, a.table))
	var amountAccounts int
	err := row.Scan(&amountAccounts)
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(context.Background(), fmt.Sprintf(`SELECT * FROM %s`, a.table))
	if err != nil {
		return nil, err
	}

	accounts := make([]account.Account, amountAccounts)

	var (
		id, gameId, ownerId, energyLimit int
		friendlyName, url                string
	)

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&id, &gameId, &friendlyName, &ownerId, &url, &energyLimit)
		if err != nil {
			return nil, err
		}
		accounts[i] = account.Account{
			ID:           id,
			GameID:       gameId,
			FriendlyName: friendlyName,
			Owner:        ownerId,
			URL:          url,
			EnergyLimit:  energyLimit,
		}
	}
	return accounts, nil
}

func (a *AccountStorage) GetById(id int) (account.Account, error) {
	var (
		gameId, ownerId, energyLimit int
		friendlyName, url            string
	)
	row := a.db.QueryRow(context.Background(), fmt.Sprintf(`SELECT game_id, friendly_name, owner_id, url, energy_limit FROM %s WHERE id == %d`, a.table, id))
	err := row.Scan(&gameId, &friendlyName, &ownerId, &url, &energyLimit)
	if err != nil {
		return account.Account{}, err
	}
	return account.Account{
		ID:           id,
		GameID:       gameId,
		FriendlyName: friendlyName,
		Owner:        ownerId,
		URL:          url,
	}, nil
}

func (a *AccountStorage) GetTable() string {
	return a.table
}

func (a *AccountStorage) Close() error {
	return a.db.Close(context.Background())
}
