package lru

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var cl = NewList()


func TestNewList(t *testing.T) {
	c := NewList(100)
	assert.Equal(t, 100, c.maxLength)
	assert.Equal(t, 10, cl.maxLength)
}

func TestCacheList_AddCache(t *testing.T) {
	for i := 0; i < 100; i++ {
		cl.AddCache(strconv.Itoa(i), strconv.Itoa(i))
		assert.Equal(t, len(cl.m), cl.l.Len(), "err")
		assert.LessOrEqual(t, len(cl.m), cl.maxLength, "err")
		assert.LessOrEqual(t, cl.l.Len(), cl.maxLength, "err")
	}
}

func TestCacheList_RemoveCache(t *testing.T) {
	for i := 0; i < 10; i++ {
		cl.AddCache(strconv.Itoa(i), strconv.Itoa(i))
	}
	cl.RemoveCache("1")
	a, ok := cl.GetCache("1")
	assert.Equal(t, a, nil, "err")
	assert.Equal(t, ok, false, "err")
}

func TestCacheList_GetCache(t *testing.T) {
	for i := 0; i < 100; i++ {
		cl.AddCache(strconv.Itoa(i), strconv.Itoa(i))
	}
	a1, ok1 := cl.GetCache("95")
	a2, ok2 := cl.GetCache("5")
	assert.Equal(t, strconv.Itoa(95), a1, "err")
	assert.Equal(t, cl.m[cl.l.Front().Value.(string)], a1, "err")
	assert.True(t, ok1, "err")
	assert.Equal(t, nil, a2, "err")
	assert.False(t, ok2, "err")
}
