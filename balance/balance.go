package balance

import "net/http"

type LoadBalance interface {
	Next(r *http.Request) *Node // 获取下一个节点
}

type Node struct {
	name          string
	weight        int
	currentWeight int
}