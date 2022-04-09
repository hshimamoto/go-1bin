// 1bin / gorkscrew
// MIT License Copyright(c) 2022 Hiroshi Shimamoto
// vim:set sw=4 sts=4:

package gorkscrew

import (
    "log"
    "net"
    "syscall"
)

var usage string = "gorkscrew [-F] <proxyhost:port> <remotehost:port>"

func passfd(conn net.Conn) {
    var raw syscall.RawConn
    var err error
    switch sock := conn.(type) {
    case *net.TCPConn:
	raw, err = sock.SyscallConn()
    case *net.UnixConn:
	raw, err = sock.SyscallConn()
    default:
	log.Println("unspported type")
	return
    }
    if err != nil {
	log.Printf("unable to get RawConn: %v\n", err)
	return
    }
    raw.Control(func(fd uintptr) {
	rights := syscall.UnixRights(int(fd))
	err := syscall.Sendmsg(1, nil, rights, nil, 0)
	if err != nil {
	    log.Printf("unable to send rights: %v\n", err)
	}
    })
}
