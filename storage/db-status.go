package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xaosBotTeam/go-shared-models/status"
)

type AbstractStatusStorage interface {
	GetById(id int) (status.Status, error)
	GetAll() (map[int]status.Status, error)
	Update(id int, stat status.Status) error
	Add(id int, stat status.Status) error
	Close()
	Delete(id int) error
}

func NewStatusStorage(connStr string) (*StatusStorage, error) {
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	table := "bot.Status"
	createSchemaString := `CREATE SCHEMA IF NOT EXISTS bot`
	createTableString := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id  		  int PRIMARY KEY,
	game_id      int NOT NULL,
	friendly_name           text NOT NULL,
	energy_limit int NOT NULL)`, table)

	_, err = conn.Exec(context.Background(), createSchemaString)
	if err != nil {
		return nil, err
	}
	_, err = conn.Exec(context.Background(), createTableString)
	if err != nil {
		return nil, err
	}

	return &StatusStorage{db: conn,
		table: table}, nil
}

type StatusStorage struct {
	db    *pgxpool.Pool
	table string
}

func (s *StatusStorage) GetById(id int) (status.Status, error) {
	var (
		gameId, energyLimit int
		friendlyName        string
	)
	row := s.db.QueryRow(context.Background(), fmt.Sprintf(`SELECT game_id, friendly_name, energy_limit FROM %s WHERE id = %d`, s.table, id))
	err := row.Scan(&gameId, &friendlyName, &energyLimit)
	if err != nil {
		return status.Status{}, err
	}
	return status.Status{
		GameID:       gameId,
		FriendlyName: friendlyName,
		EnergyLimit:  energyLimit,
	}, nil
}

func (s *StatusStorage) GetAll() (map[int]status.Status, error) {
	rows, err := s.db.Query(context.Background(), fmt.Sprintf(`SELECT * FROM %s`, s.table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	statuses := make(map[int]status.Status)

	var (
		id, gameId, energyLimit int
		friendlyName            string
	)

	for rows.Next() {
		err = rows.Scan(&id, &gameId, &friendlyName, &energyLimit)
		if err != nil {
			return nil, err
		}
		statuses[id] = status.Status{
			GameID:       gameId,
			FriendlyName: friendlyName,
			EnergyLimit:  energyLimit,
		}
	}
	return statuses, nil
}

func (s *StatusStorage) Update(id int, stat status.Status) error {
	_, err := s.db.Exec(context.Background(), fmt.Sprintf("UPDATE %s SET (game_id, friendly_name, energy_limit)" +
		" = (%d, '%s', %d) WHERE id = %d", s.table, stat.GameID, stat.FriendlyName, stat.EnergyLimit, id))
	return err
}

func (s *StatusStorage) Add(id int, stat status.Status) error {
	_, err := s.db.Exec(context.Background(), fmt.Sprintf("INSERT INTO %s VALUES (%d, %d, '%s', %d)", s.table, id, stat.GameID, stat.FriendlyName, stat.EnergyLimit))
	return err
}

func (s *StatusStorage) Close() {
	s.db.Close()
}

func (s *StatusStorage) Delete(id int) error {
	_, err := s.db.Exec(context.Background(), fmt.Sprintf(`DELETE FROM %s WHERE id = %d`, s.table, id))
	return err
}
