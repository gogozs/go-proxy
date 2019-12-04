package main

import (
	"fmt"
	"go-proxy/conf"
	"go-proxy/route"
	"log"
	"net/http"
	_ "net/http/pprof"  // pprof debug
	"path"
	"time"
)

func startServer() {
	c := conf.GetConfig()
	r := route.GetRouter()
	readTimeOut := c.Common.ReadTimeOut
	writeTimeOut := c.Common.WriteTimeOut
	if readTimeOut <= 0 {
		readTimeOut = 5
	}
	if writeTimeOut <= 0 {
		readTimeOut = 10
	}
	server := http.Server{
		ReadTimeout: time.Duration(readTimeOut) * time.Second,
		WriteTimeout: time.Duration(writeTimeOut) * time.Second,
		Handler: r,
		Addr: fmt.Sprintf(":%s", c.Server.Port),
	}
	if c.Common.Tls {
		// https
		confPath := conf.GetConfigPath()
		crtPath := path.Join(confPath, "server.crt")
		keyPath := path.Join(confPath, "server.key")
		if e := server.ListenAndServeTLS(crtPath, keyPath); e != nil {
			log.Fatalln("ListenAndServe: ", e)
		}
	} else {
		log.Println(fmt.Sprintf("程序运行：http://localhost:%s", c.Server.Port))
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}
}

func main() {
	//go http.ListenAndServe(":8080", nil) // debug
	startServer()
}