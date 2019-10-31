package conf

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitVal(t *testing.T) {
	vc := VarConfig(map[string]string{
		"GOROOT": "$GOROOT",
	})
	c := &Config{Var:vc}
	initVars(c)
	target := os.Getenv("GOROOT")
	assert.Equal(t, c.Var["GOROOT"], target)
}
