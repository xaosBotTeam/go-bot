package storage

import (
	"errors"
	models "github.com/xaosBotTeam/go-shared-models/task"
)

func NewFakeStatusStorage() AbstractStatusStorage {
	return &FakeStatusStorage{statuses: make(map[int]models.Status)}
}

type FakeStatusStorage struct {
	statuses map[int]models.Status
}

func (f *FakeStatusStorage) GetAll() ([]int, []models.Status, error) {
	ids := make([]int, len(f.statuses))
	statuses := make([]models.Status, len(f.statuses))
	i := 0
	for id, status := range f.statuses {
		ids[i] = id
		statuses[i] = status
		i++
	}
	return ids, statuses, nil
}

func (f *FakeStatusStorage) GetByAccId(id int) (models.Status, error) {
	if status, ok := f.statuses[id]; ok {
		return status, nil
	}
	return models.Status{}, errors.New("not found")
}

func (f *FakeStatusStorage) Update(id int, status models.Status) error {
	f.statuses[id] = status
	return nil
}

func (f *FakeStatusStorage) Delete(id int) error {
	delete(f.statuses, id)
	return nil
}

func (f *FakeStatusStorage) Add(id int, status models.Status) error {
	f.statuses[id] = status
	return nil
}
