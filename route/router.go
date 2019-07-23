package route

import (
	"fmt"
	"go-proxy/conf"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type router struct {
	rules []conf.ProxyConfig
}

var r = &router{}

func init() {
	r.rules = conf.GetConfig().Proxy
	fmt.Println(r.rules[0])
}

func (this *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	for _, proxy := range this.rules {
		re, _ := regexp.Compile(proxy.Location)
		if re.MatchString(path) {
			r.URL.Path = re.ReplaceAllString(path, "/") // 代理
			remote, err := url.Parse(proxy.ProxyPass)
			if err != nil {
				log.Println(err)
				panic(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(remote)
			proxy.ServeHTTP(w, r)
			break
		}
	}
	// TODO 跳转404  根目录指向前端静态文件

}

func GetRouter() *router {
	return r
}
