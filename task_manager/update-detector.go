package task_manager

import "sync"

func NewUpdateDetector() *UpdateDetector {
	return &UpdateDetector{
		mx: sync.RWMutex{},
		m:  make(map[int]bool),
	}
}

type UpdateDetector struct {
	mx sync.RWMutex
	m  map[int]bool
}

func (u *UpdateDetector) Load(key int) bool {
	u.mx.RLock()
	val := u.m[key]
	u.mx.RUnlock()
	return val
}

func (u *UpdateDetector) Store(key int, val bool) {
	u.mx.Lock()
	u.m[key] = val
	u.mx.Unlock()
}

func (u *UpdateDetector) Delete(key int) {
	u.mx.Lock()
	delete(u.m, key)
	u.mx.Unlock()
}
