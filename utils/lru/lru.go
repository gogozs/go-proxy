/*
	simple lru
 */

package lru

import "container/list"

type Cache struct {
	key   string
	value string
}

type CacheList struct {
	l         *list.List
	m         map[string]string
	maxLength int
}

func NewList(num ...int) *CacheList {
	maxLength := 10
	if len(num) > 0 && num[0] > 0 {
		maxLength = num[0]
	}
	return &CacheList{list.New(), make(map[string]string), maxLength}
}

func (cl *CacheList) GetCache(key string) (string, bool) {
	if value, ok := cl.m[key]; ok {
		cl.MoveFront(Cache{key:key, value:value})
		return value, ok
	} else {
		return "", false
	}
}

func (cl *CacheList) AddCache(key, value string) {
	c := Cache{key: key, value: value}
	if _, ok := cl.m[c.key]; !ok {
		cl.m[c.key] = c.value
		cl.l.PushFront(c.key)
		if cl.l.Len() > cl.maxLength {
			delete(cl.m, cl.l.Back().Value.(string))
			cl.l.Remove(cl.l.Back())
		}
	} else {
		cl.m[c.key] = c.value
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
	delete(cl.m, key)
	for e := cl.l.Front(); e != nil; e = e.Next() {
		if e.Value == key {
			cl.l.Remove(e)
		}
	}
}
