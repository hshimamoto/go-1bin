// 1bin / gorkscrew
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
// +build linux

package gorkscrew

import (
    "io"
    "log"
    "net"
    "os"
    "syscall"

    "github.com/hshimamoto/go-iorelay"
    "github.com/hshimamoto/go-session"
)

type stdio struct {
    stdin io.Reader
    stdout io.Writer
}

func (s *stdio)Read(p []byte) (int, error) {
    return s.stdin.Read(p)
}

func (s *stdio)Write(p []byte) (int, error) {
    return s.stdout.Write(p)
}

func Run(args []string) {
    // gorkscrew proxyhost:port remotehost:port
    if len(args) < 2 {
	log.Println("gorkscrew [-F] <proxyhost:port> <remotehost:port>")
	return
    }
    opt := ""
    idx := 0
    if len(args) >= 3 {
	opt = args[0]
	idx = 1
    }
    proxyaddr := args[idx]
    remoteaddr := args[idx + 1]
    if opt != "" && opt != "-F" {
	log.Printf("unknown opt: %v", args)
	return
    }

    // Do HTTP CONNECT Here
    conn, err := session.Corkscrew(proxyaddr, remoteaddr)
    if err != nil {
	log.Printf("Corkscrew %s %s error: %v\n", proxyaddr, remoteaddr, err)
	return
    }

    if opt == "-F" {
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
	// finish passfd
	return
    }

    iorelay.Relay(conn, &stdio{ stdin: os.Stdin, stdout: os.Stdout })
}
