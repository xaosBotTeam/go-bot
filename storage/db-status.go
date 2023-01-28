package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	models "github.com/xaosBotTeam/go-shared-models/task"
)

type AbstractStatusStorage interface {
	GetAll() ([]int, []models.Status, error)
	GetByAccId(id int) (models.Status, error)
	Update(id int, status models.Status) error
	Delete(id int) error
	Add(id int, status models.Status) error
}

func NewStatusStorage(connString string) (AbstractStatusStorage, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	table := "Status"
	createSchemaString := `CREATE SCHEMA IF NOT EXISTS bot`
	createTableString := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id int PRIMARY KEY,
	status json)`, table)

	_, err = conn.Exec(context.Background(), createSchemaString)
	if err != nil {
		return nil, err
	}
	_, err = conn.Exec(context.Background(), createTableString)
	if err != nil {
		return nil, err
	}

	return &DbStatusStorage{db: conn,
		table: table}, nil
}

type DbStatusStorage struct {
	db    *pgx.Conn
	table string
}

func (d *DbStatusStorage) GetAll() ([]int, []models.Status, error) {
	row := d.db.QueryRow(context.Background(), fmt.Sprintf(`SELECT COUNT(*) FROM %s`), d.table)
	var amountRows int
	err := row.Scan(&amountRows)
	if err != nil {
		return nil, nil, err
	}

	rows, err := d.db.Query(context.Background(), fmt.Sprint("SELECT * FROM %s", d.table))
	if err != nil {
		return nil, nil, err
	}

	statuses := make([]models.Status, amountRows)
	ids := make([]int, amountRows)
	var (
		id            int
		statusJsonStr string
		status        models.Status
	)

	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&id, &statusJsonStr)
		if err != nil {
			return nil, nil, err
		}
		err = json.Unmarshal([]byte(statusJsonStr), &status)
		if err != nil {
			return nil, nil, err
		}

		statuses[i] = status
		ids[i] = id
	}
	return ids, statuses, nil
}

func (d *DbStatusStorage) GetByAccId(id int) (models.Status, error) {
	var (
		statusJsonStr string
		status        models.Status
	)

	row := d.db.QueryRow(context.Background(), fmt.Sprintf(`SELECT status FROM %s WHERE id == %d`), d.table, id)
	err := row.Scan(&statusJsonStr)
	if err != nil {
		return models.Status{}, err
	}

	err = json.Unmarshal([]byte(statusJsonStr), &status)
	if err != nil {
		return models.Status{}, err
	}

	return status, nil
}

func (d *DbStatusStorage) Update(id int, status models.Status) error {
	jsonStr, err := json.Marshal(status)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(context.Background(), fmt.Sprintf("UPDATE %s SET status = '%s' WHERE id = %d", d.table, id, jsonStr))
	return err
}

func (d *DbStatusStorage) Delete(id int) error {
	_, err := d.db.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s WHERE id = %d", d.table, id))
	return err
}

func (d *DbStatusStorage) Add(id int, status models.Status) error {
	jsonStr, err := json.Marshal(status)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(context.Background(), fmt.Sprintf("INSERT INTO %s VALUES (%d, '%s')", d.table, id, string(jsonStr)))
	return err
}
