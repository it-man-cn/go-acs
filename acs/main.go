package main

import (
	"github.com/astaxie/beego"
	_ "go-acs/acs/routers"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
	beego.Run()
}
