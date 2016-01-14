package main

import (
	log "code.google.com/p/log4go"
	"fmt"
	"net"
)

func main() {
	var (
		addr     *net.UDPAddr
		listener *net.UDPConn
		err      error
	)
	log.LoadConfiguration(Conf.Log)
	defer log.Close()
	if addr, err = net.ResolveUDPAddr("udp4", Conf.Bind); err != nil {
		log.Error("net.ResolveUDPAddr(\"udp4\", \"%s\") error(%v)", Conf.Bind, err)
		return
	}

	if listener, err = net.ListenUDP("udp4", addr); err != nil {
		log.Error("net.ListenUDP(\"udp4\", \"%v\") error(%v)", addr, err)
		return
	}
	defer listener.Close()

	if Debug {
		log.Debug("start udp listen: \"%s\"", Conf.Bind)
	}

	//N core accept
	for i := 0; i < Conf.MaxProc; i++ {
		go acceptUDP(listener)
	}
	//wait
	InitSignal()
}

func acceptUDP(listener *net.UDPConn) {
	for {
		buff := make([]byte, 1024)
		n, addr, err := listener.ReadFromUDP(buff)
		if err != nil {
			log.Error("ReadFromUdp(\"%v\") error(%v) nbytes(%d)", addr, err, n)
		}
		if n > 0 {
			fmt.Printf("read %s", string(buff))
		}

	}
}
