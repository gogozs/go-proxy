package balance

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ih *IpHashmap

func init() {
	n1 := &Node{name: "n1"}
	n2 := &Node{name: "n2"}
	n3 := &Node{name: "n3"}
	n4 := &Node{name: "n4"}
	ih = &IpHashmap{nodes:[]*Node{n1, n2, n3, n4}}
}


func TestIpHashmap_Next(t *testing.T) {
	ips := []string{
		"10.0.1.0:8000",
		"10.0.1.0:9000",
		"100.1.1.0:9000",
		"101.1.4.0:9000",
		"105.1.1.0:9000",
		"205.3.5.0:9000",
		"205.4.5.0:9000",
		"205.5.5.0:9000",
		"205.6.5.0:9000",
		"205.7.5.0:9000",
	}
	for _, ip := range ips {
		h := ih.getHash(ip)
		fmt.Println(h)
		assert.Less(t, h, len(ih.nodes))
		assert.GreaterOrEqual(t, h, 0)
	}
}