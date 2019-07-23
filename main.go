package main

import (
	"fmt"
	"go-proxy/conf"
	"go-proxy/route"
	"log"
	"net/http"
)

func startServer() {
	c := conf.GetConfig().Server
	r := route.GetRouter()
	log.Println(fmt.Sprintf("程序运行：http://localhost:%s", c.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%s", c.Port), r)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

func main() {
	startServer()
}