package task_manager

import (
	"time"
)

type AbstractAccountQueue interface {
	Add(id int)
	Get() int
}

func NewAccountQueue() *AccountQueue {
	return &AccountQueue{ids: make([]int, 0)}
}

type AccountQueue struct {
	ids []int
}


func (a *AccountQueue) Add(id int) {
	a.ids = append(a.ids, id)
}

func (a *AccountQueue) Get() int {
	for {
		if len(a.ids) > 0 {
			value := a.ids[0]
			a.ids = a.ids[1:]
			return value
		}
		time.Sleep(100 * time.Millisecond)
	}
}
