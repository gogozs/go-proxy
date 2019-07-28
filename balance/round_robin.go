package balance

import "net/http"

/*
轮询算法
Round Robin Scheduling - 轮询调度算法
*/
type NormalRound struct {
	index int // 当前节点号
	nodes []*Node
}

func (this *NormalRound) Next(r *http.Request) *Node {
	n := this.nodes[this.index]
	if this.index == len(this.nodes) - 1 {
		this.index = 0
	} else {
		this.index++
	}
	return n
}


// 参考nginx, 平滑的加权轮询（smooth weighted round-robin balancing）的算法
type WeightRound struct {
	total int // 节点权重之和
	nodes []*Node
}

func (this *WeightRound) Next(r *http.Request) *Node {
	this.Round()
	var max int
	var node *Node
	for _, n := range this.nodes {
		if n.currentWeight > max {
			max = n.currentWeight
			node = n
		}
	}
	node.currentWeight -= this.total // 选中之后需减少总权重
	return node
}

// 回合
func (this *WeightRound) Round() {
	for _, n := range this.nodes {
		n.currentWeight += n.weight // 每回合增加自身权重
	}
}


// 工厂
func RoundBuilder(roundType string, nodes []*Node) LoadBalance {
	if roundType == "weight" {
		var total int
		for _, n := range nodes {
			total += n.weight
		}
		return &WeightRound{
			total: total,
			nodes: nodes,
		}
	} else {
		return &NormalRound{
			index: 0,
			nodes: nodes,
		}
	}
}
