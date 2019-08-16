package route

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRouter_ServeHTTP(t *testing.T) {
	c := http.Client{}
	req1, err1 := http.NewRequest("GET", "http://localhost/", nil)
	req2, err2 := http.NewRequest("GET", "http://localhost/index.html", nil)
	req3, err3 := http.NewRequest("GET", "http://localhost/test.zip", nil)
	res1, err4 := c.Do(req1)
	res2, err5 := c.Do(req2)
	res3, err6 := c.Do(req3)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
	assert.Nil(t, err4)
	assert.Nil(t, err5)
	assert.Nil(t, err6)
	assert.Equal(t, 200, res1.StatusCode)
	assert.Equal(t, 200, res2.StatusCode)
	assert.Equal(t, 200, res3.StatusCode)
	r1, _ := ioutil.ReadAll(res1.Body)
	r2, _ := ioutil.ReadAll(res2.Body)
	r3, _ := ioutil.ReadAll(res3.Body)
	assert.Equal(t, r1, r2)
	assert.NotEqual(t, r2, r3)
}

func TestExists(t *testing.T) {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)

	assert.Equal(t, true, Exists("/etc"))
	assert.Equal(t, false, Exists("/etc/stxdfd/"))
	assert.Equal(t, true, Exists(path))
	assert.Equal(t, false, Exists("/etc/1sdf.tt"))

}

func TestGetRouter(t *testing.T) {
	r := GetRouter()
	assert.NotNil(t, r)
}