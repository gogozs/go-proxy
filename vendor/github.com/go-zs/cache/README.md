## a simple golang cache store
很多简单服务需要用到缓存，但是又不需要单独引入redis。
使用标准库 sync.Map作为缓存存储，可以做到高并发下线程安全。

## Usage example

```golang

package main

import (
	"fmt"
	"github.com/go-zs/cache"
)

func main() {
	store := cache.NewList()
	store.SetCache("t", 15)
	t, _ := store.GetCache("t")
	fmt.Println(t)
}

```
