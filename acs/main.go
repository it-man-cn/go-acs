package main

import (
	//"github.com/astaxie/beego"
	//_ "go-acs/acs/routers"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//beego.Run()

	addrs := []string{":10090"}
	initHTTP(addrs)

	// block until a signal is received.
	InitSignal()

}

// InitSignal register signals handler.
func InitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		fmt.Printf(" get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			return
		case syscall.SIGHUP:
			//reload()
		default:
			return
		}
	}
}
