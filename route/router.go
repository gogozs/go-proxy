package route

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-zs/cache"
	"go-proxy/conf"
	"go-proxy/log"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type router struct {
	rules       []conf.ProxyConfig
	staticPaths []conf.StaticConfig
}

const basePath = "./html"

var (
	r         = &router{}
	cachePath = cache.NewStore()
)

func init() {
	r.rules = conf.GetConfig().Server.Proxy
	r.staticPaths = conf.GetConfig().Server.Static
}

// 判断静态文件是否存在
func checkStaticfile(urlPath, basePath string) bool {
	file := path.Join(basePath, urlPath)
	return Exists(file)
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (this *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	log.Info(fmt.Sprintf("request: %s", urlPath))
	if urlPath == "/" || checkStaticfile(urlPath, basePath) {
		this.ServeStatic(w, r, basePath) // 根目录指向前端静态文件
		return
	}

	for _, staticConf := range this.staticPaths {
		if strings.HasPrefix(urlPath, staticConf.Path) {
			r.URL.Path = strings.TrimPrefix(urlPath, staticConf.Path)
			this.ServeStatic(w, r, staticConf.Alias)
			return
		}
	}

	if proxyPass, ok := cachePath.GetCache(urlPath); ok {
		this.ServeProxy(w, r, proxyPass.(string))
		return
	} else {
		for _, proxy := range this.rules {
			re, _ := regexp.Compile(proxy.Location)
			if re.MatchString(urlPath) {
				proxyPass := proxy.ProxyPass
				log.Info(fmt.Sprintf("%s: %s", urlPath, proxy.ProxyPass))
				this.ServeProxy(w, r, proxyPass)
				cachePath.SetCache(urlPath, proxyPass)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 not found")
}

func (this *router) ServeProxy(w http.ResponseWriter, r *http.Request, proxyPass string) {
	remote, err := url.Parse(proxyPass)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	// 设置超时时间
	t := conf.GetConfig().Common.Timeout
	if t < 5 {
		t = 15 // 默认15秒超时
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second * time.Duration(t))
	r = r.WithContext(ctx)
	proxy.ServeHTTP(w, r)
}

func newCustomerHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}

	return &httputil.ReverseProxy{
		Director: director,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, time.Second*30)},
		},
	}
}


// 代理静态文件
func (this *router) ServeStatic(w http.ResponseWriter, r *http.Request, path string) {
	handler := http.FileServer(http.Dir(path))
	handler.ServeHTTP(w, r)
}

// 提供静态文件下载
func (this *router) ServeDownload(w http.ResponseWriter, r *http.Request, filePath, fileName string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprint(w, err)
	}
	http.ServeContent(w, r, fileName, time.Now(), bytes.NewReader(data))
}

func GetRouter() *router {
	return r
}
