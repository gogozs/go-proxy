package route

import (
	"bytes"
	"fmt"
	"go-proxy/conf"
	"go-proxy/log"
	"go-proxy/utils/lru"
	"io/ioutil"
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
	cachePath = lru.NewList()
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
			r.URL.Path = strings.TrimLeft(urlPath, staticConf.Path)
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
				log.Info(fmt.Sprintf("%s: %s", urlPath,  proxy.ProxyPass))
				this.ServeProxy(w, r, proxyPass)
				cachePath.SetCache(urlPath,proxyPass)
				return
			}
		}
	}
}

func (this *router) ServeProxy(w http.ResponseWriter, r *http.Request, proxyPass string) {
	remote, err := url.Parse(proxyPass)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

// 代理静态文件
func (this *router) ServeStatic(w http.ResponseWriter, r *http.Request, path string) {
	hander := http.FileServer(http.Dir(path))
	hander.ServeHTTP(w, r)
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
