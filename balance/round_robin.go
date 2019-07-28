package balance

/*
轮询算法
Round Robin Scheduling - 轮询调度算法
*/

type RoundRobin interface {
	Next() *Node // 获取下一个节点
}

type Node struct {
	name          string
	weight        int
	currentWeight int
}


// 参考nginx, 平滑的加权轮询（smooth weighted round-robin balancing）的算法
type NginxRound struct {
	total int // 节点权重之和
	nodes []*Node
}

func (this *NginxRound) Next() *Node {
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
func (this *NginxRound) Round() {
	for _, n := range this.nodes {
		n.currentWeight += n.weight // 每回合增加自身权重
	}
}


// 工厂
func RoundBuilder(roundType string, nodes []*Node) RoundRobin {
	var total int
	for _, n := range nodes {
		total += n.weight
	}
	return &NginxRound{
		total: total,
		nodes: nodes,
	}
}
