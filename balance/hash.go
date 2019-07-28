package balance

import (
	"hash/fnv"
	"net/http"
)

// 根据host hash值选择对应的a节点
type IpHashmap struct {
	nodes []*Node
}

func (this *IpHashmap) Next(r *http.Request) *Node{
	host := r.RemoteAddr  // host = ip+port
	hh := this.getHash(host)
	return this.nodes[hh]
}

// 求hash值
func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func (this *IpHashmap) getHash(host string) int {
	h := hash(host)
	hh := h % len(this.nodes)
	return hh
}