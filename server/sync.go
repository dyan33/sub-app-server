package server

import (
	"sync"
)

type ID struct {
	id int64

	mutex *sync.Mutex
}

func (i *ID) get() int64 {

	i.mutex.Lock()
	defer func() { i.mutex.Unlock() }()

	i.id = i.id + 1

	return i.id
}

func (i *ID) set(value int64) {
	i.mutex.Lock()
	defer func() { i.mutex.Unlock() }()
	i.id = value
}

func newID() *ID {

	return &ID{
		id:    0,
		mutex: &sync.Mutex{},
	}
}
