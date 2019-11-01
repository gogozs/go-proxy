package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func StartServer() {
	http.HandleFunc("/", handler)
	crtPath := "./conf/server.crt"
	keyPath := "./conf/server.key"
	if e := http.ListenAndServeTLS(fmt.Sprintf(":%s", "10086"), crtPath, keyPath, nil); e != nil {
		log.Fatal("ListenAndServe: ", e)
	}
}

func main() {
	//go http.ListenAndServe(":8888", nil) # debug
	StartServer()
}