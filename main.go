package main

import (
	"fmt"
	"go-proxy/conf"
	"go-proxy/route"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path"
)

func startServer() {
	c := conf.GetConfig()
	r := route.GetRouter()
	if c.Common.Tls {
		// https
		confPath := conf.GetConfigPath()
		crtPath := path.Join(confPath, "server.crt")
		keyPath := path.Join(confPath, "server.key")
		if e := http.ListenAndServeTLS(fmt.Sprintf(":%s", c.Server.Port), crtPath, keyPath, nil); e != nil {
			log.Fatal("ListenAndServe: ", e)
		}
	} else {
		log.Println(fmt.Sprintf("程序运行：http://localhost:%s", c.Server.Port))
		err := http.ListenAndServe(fmt.Sprintf(":%s", c.Server.Port), r)
		if err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}
}

func main() {
	//go http.ListenAndServe(":8888", nil) # debug
	startServer()
}