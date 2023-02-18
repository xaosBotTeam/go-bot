package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xaosBotTeam/go-shared-models/account"
)

func NewAccountStorage(connString string) (*AccountStorage, error) {
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
	owner_id      int NOT NULL,
	url           text NOT NULL)`, table)

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

type AccountStorage struct {
	db    *pgxpool.Pool
	table string
}

func (a *AccountStorage) GetAll() (map[int]account.Account, error) {
	rows, err := a.db.Query(context.Background(), fmt.Sprintf(`SELECT * FROM %s`, a.table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	accounts := make(map[int]account.Account)

	var (
		id, ownerId int
		url         string
	)

	for rows.Next() {
		err = rows.Scan(&id, &ownerId, &url)
		if err != nil {
			return nil, err
		}
		accounts[id] = account.Account{
			Owner: ownerId,
			URL:   url,
		}
	}
	return accounts, nil
}

func (a *AccountStorage) GetById(id int) (account.Account, error) {
	var (
		ownerId int
		url            string
	)
	row := a.db.QueryRow(context.Background(), fmt.Sprintf(`SELECT owner_id, url FROM %s WHERE id = %d`, a.table, id))
	err := row.Scan(&ownerId, &url)
	if err != nil {
		return account.Account{}, err
	}
	return account.Account{
		Owner: ownerId,
		URL:   url,
	}, nil
}

func (a *AccountStorage) Close() {
	a.db.Close()
}

func (a *AccountStorage) Add(acc account.Account) (int, error) {
	var id int
	query := a.db.QueryRow(context.Background(), fmt.Sprintf(`INSERT INTO %s (owner_id, url) VALUES (%d, '%s')
                               RETURNING id`, a.table, acc.Owner, acc.URL))
	err := query.Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (a *AccountStorage) Update(id int, acc account.Account) error {
	_, err := a.db.Exec(context.Background(), fmt.Sprintf(`UPDATE %s SET (owner_id, url) = (%d, '%s') WHERE id = %d`, a.table, acc.Owner, acc.URL, id))
	return err
}

func (a *AccountStorage) Delete(id int) error {
	_, err := a.db.Exec(context.Background(), fmt.Sprintf(`DELETE FROM %s WHERE id = %d`, a.table, id))
	return err
}
