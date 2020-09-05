// 1bin / fwd
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package fwd

import (
    "log"
    "net"

    "github.com/hshimamoto/go-session"
    "github.com/hshimamoto/go-iorelay"
)

func Run(args []string) {
    if len(args) < 2 {
	log.Println("fwd <listen> <dst>")
	return
    }
    serv, err := session.NewServer(args[0], func(conn net.Conn) {
	defer conn.Close()
	fconn, err := session.Dial(args[1])
	if err != nil {
	    log.Printf("Dial %s %v\n", args[1], err)
	    return
	}
	defer fconn.Close()
	iorelay.Relay(conn, fconn)
    })
    if err != nil {
	log.Printf("NewServer %s %v\n", args[0], err)
	return
    }
    serv.Run()
}
