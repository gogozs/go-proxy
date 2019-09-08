/*
	simple lru
 */

package lru

import (
	"container/list"
	"sync"
)

type Cache struct {
	key   string
	value interface{}
}

type CacheList struct {
	l         *list.List
	m         sync.Map
	maxLength int
}

func NewList(num ...int) *CacheList {
	maxLength := 1000
	if len(num) > 0 && num[0] > 0 {
		maxLength = num[0]
	}
	return &CacheList{list.New(), sync.Map{}, maxLength}
}

func (cl *CacheList) GetCache(key string) (interface{}, bool) {
	if value, ok := cl.m.Load(key); ok {
		cl.MoveFront(Cache{key:key, value:value})
		return value, ok
	} else {
		return "", false
	}
}

func (cl *CacheList) SetCache(key string, value interface{}) {
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

func (cl *CacheList) MoveFront(c Cache) {
	for e := cl.l.Front(); e != nil; e = e.Next() {
		if e.Value == c.key {
			cl.l.MoveToFront(e)
		}
	}
}

func (cl *CacheList) RemoveCache(key string) {
	cl.m.Delete(key)
	for e := cl.l.Front(); e != nil; e = e.Next() {
		if e.Value == key {
			cl.l.Remove(e)
		}
	}
}
