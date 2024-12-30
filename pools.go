package core

import "sync"

type Pool[T any] struct {
	pool *sync.Pool
}

func NewPool[T any]() *Pool[T] {
	return &Pool[T]{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(T)
			},
		},
	}
}

func (that *Pool[T]) Put(entry *T) {
	that.pool.Put(entry)
}

func (that *Pool[T]) Get() *T {
	return that.pool.Get().(*T)
}
