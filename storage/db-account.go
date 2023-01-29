package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xaosBotTeam/go-shared-models/account"
)

func NewAccountStorage(connString string) (AbstractAccountStorage, error) {
	conn, err := pgxpool.New(context.Background(), connString)
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
    id  		  SERIAL PRIMARY KEY,
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
	Close()
	Add(url string, ownerId int) (account.Account, error)
	Update(acc account.Account) error
}

type AccountStorage struct {
	db    *pgxpool.Pool
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
	defer rows.Close()
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
	row := a.db.QueryRow(context.Background(), fmt.Sprintf(`SELECT game_id, friendly_name, owner_id, url, energy_limit FROM %s WHERE id = %d`, a.table, id))
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

func (a *AccountStorage) Close() {
	a.db.Close()
}

func (a *AccountStorage) Add(url string, ownerId int) (account.Account, error) {
	var id int
	query := a.db.QueryRow(context.Background(), fmt.Sprintf(`INSERT INTO %s (game_id, friendly_name, owner_id, url, energy_limit)
	 VALUES (%d, '%s', %d, '%s', %d) RETURNING id`, a.table, 0, "New account", ownerId, url, 1000))
	err := query.Scan(&id)
	if err != nil {
		return account.Account{}, nil
	}
	return account.Account{
		ID:           id,
		GameID:       0,
		FriendlyName: "New account",
		Owner:        ownerId,
		URL:          url,
		EnergyLimit:  1000,
	}, nil
}

func (a *AccountStorage) Update(acc account.Account) error {
	_, err := a.db.Exec(context.Background(), fmt.Sprintf(`UPDATE %s SET (game_id, friendly_name, owner_id, url, energy_limit) =
    (%d, '%s',%d, '%s', %d) WHERE id = %d`, a.table, acc.GameID, acc.FriendlyName, acc.Owner, acc.URL, acc.EnergyLimit, acc.ID))
	return err
}
