package conf

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitVal(t *testing.T) {
	vc := VarConfig(map[string]string{
		"GOROOT": "$GOROOT",
		"RPA_APP_FILE_PATH": "/tmp",
	})
	c := &Config{Var:vc, Server:ServerConfig{Static: []StaticConfig{{Path: "/v1/path", Alias:"$RPA_APP_FILE_PATH"}}}}
	initVars(c)
	target := os.Getenv("GOROOT")
	assert.Equal(t, c.Var["GOROOT"], target)
	updateConfig(c)
	//assert.Equal(t, true, path.IsAbs(c.))
	fmt.Println(c.Server.Static)
}
