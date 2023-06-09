package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xaosBotTeam/go-shared-models/config"
)

func NewConfigStorage(connString string) (*DbConfigStorage, error) {
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	table := "bot.Config"
	createSchemaString := `CREATE SCHEMA IF NOT EXISTS bot`
	createTableString := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id int PRIMARY KEY,
	config json)`, table)

	_, err = conn.Exec(context.Background(), createSchemaString)
	if err != nil {
		return nil, err
	}
	_, err = conn.Exec(context.Background(), createTableString)
	if err != nil {
		return nil, err
	}

	return &DbConfigStorage{db: conn,
		table: table}, nil
}

type DbConfigStorage struct {
	db    *pgxpool.Pool
	table string
}

func (d *DbConfigStorage) GetAll() ([]int, []config.Config, error) {
	row := d.db.QueryRow(context.Background(), fmt.Sprintf(`SELECT COUNT(*) FROM %s`, d.table))
	var amountRows int
	err := row.Scan(&amountRows)
	if err != nil {
		return nil, nil, err
	}
	rows, err := d.db.Query(context.Background(), fmt.Sprintf(`SELECT * FROM %s`, d.table))
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	configs := make([]config.Config, amountRows)
	ids := make([]int, amountRows)
	var (
		id            int
		configJsonStr string
		configuration config.Config
	)

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&id, &configJsonStr)
		if err != nil {
			return nil, nil, err
		}
		err = json.Unmarshal([]byte(configJsonStr), &configuration)
		if err != nil {
			return nil, nil, err
		}

		configs[i] = configuration
		ids[i] = id
	}
	return ids, configs, nil
}

func (d *DbConfigStorage) GetByAccId(id int) (config.Config, error) {
	var (
		configJsonStr string
		configuration config.Config
	)

	row := d.db.QueryRow(context.Background(), fmt.Sprintf(`SELECT config FROM %s WHERE id = %d`, d.table, id))
	err := row.Scan(&configJsonStr)
	if err != nil {
		return config.Config{}, err
	}

	err = json.Unmarshal([]byte(configJsonStr), &configuration)
	if err != nil {
		return config.Config{}, err
	}

	return configuration, nil
}

func (d *DbConfigStorage) Update(id int, configuration config.Config) error {
	jsonStr, err := json.Marshal(configuration)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(context.Background(), fmt.Sprintf("UPDATE %s SET config = '%s' WHERE id = %d", d.table, jsonStr, id))
	return err
}

func (d *DbConfigStorage) Delete(id int) error {
	_, err := d.db.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s WHERE id = %d", d.table, id))
	return err
}

func (d *DbConfigStorage) Add(id int, configuration config.Config) error {
	jsonStr, err := json.Marshal(configuration)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(context.Background(), fmt.Sprintf("INSERT INTO %s VALUES (%d, '%s')", d.table, id, string(jsonStr)))
	return err
}

func (d *DbConfigStorage) Close() {
	d.db.Close()
}

func (d *DbConfigStorage) UpdateRange(configuration config.Config) error {
	jsonStr, err := json.Marshal(configuration)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(context.Background(), fmt.Sprintf(`UPDATE %s SET config = '%s'`, d.table, jsonStr))
	if err != nil {
		return err
	}
	return nil
}
