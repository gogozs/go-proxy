package balance

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	nr         *NginxRound
	n1, n2, n3 *Node
)

func init() {
	n1 = &Node{weight: 1, name: "n1"}
	n2 = &Node{weight: 2, name: "n2"}
	n3 = &Node{weight: 3, name: "n3"}
	nr = &NginxRound{
		total: n1.weight + n2.weight + n3.weight,
		nodes: []*Node{n1, n2, n3},
	}
}

func TestNginxRound_Next(t *testing.T) {
	var n1_count, n2_count, n3_count int
	for i := 0; i < 600; i++ {
		n := nr.Next()
		if n == n1 {
			n1_count++
		} else if n == n2 {
			n2_count++
		} else if n == n3 {
			n3_count++
		}
	}
	fmt.Println(n1_count, n2_count, n3_count)
	assert.Equal(t, n1_count, 100, "err")
	assert.Equal(t, n2_count, 200, "err")
	assert.Equal(t, n3_count, 300, "err")
}
