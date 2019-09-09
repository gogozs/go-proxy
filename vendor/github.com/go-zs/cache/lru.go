/*
	simple lru
 */

package cache

import (
	"container/list"
	"sync"
)

type Cache struct {
	key   string
	value interface{}
}

type Store struct {
	l         *list.List
	m         sync.Map
	maxLength int
}

func NewStore(num ...int) *Store {
	maxLength := 1000
	if len(num) > 0 && num[0] > 0 {
		maxLength = num[0]
	}
	return &Store{list.New(), sync.Map{}, maxLength}
}


func (cl *Store) GetCache(key string) (interface{}, bool) {
	if value, ok := cl.m.Load(key); ok {
		cl.MoveFront(Cache{key:key, value:value})
		return value, ok
	} else {
		return "", false
	}
}

func (cl *Store) SetCache(key string, value interface{}) {
	c := Cache{key: key, value: value}
	if _, ok := cl.m.Load(c.key); !ok {
		cl.m.Store(c.key, c.value)
		cl.l.PushFront(c.key)
		if cl.l.Len() > cl.maxLength {
			cl.m.Delete(cl.l.Back().Value.(string))
			cl.l.Remove(cl.l.Back())
		}
	} else {
		cl.m.Store(c.key, c.value)
		cl.MoveFront(c)
	}
}

func (cl *Store) MoveFront(c Cache) {
	for e := cl.l.Front(); e != nil; e = e.Next() {
		if e.Value == c.key {
			cl.l.MoveToFront(e)
		}
	}
}

func (cl *Store) RemoveCache(key string) {
	cl.m.Delete(key)
	for e := cl.l.Front(); e != nil; e = e.Next() {
		if e.Value == key {
			cl.l.Remove(e)
		}
	}
}
