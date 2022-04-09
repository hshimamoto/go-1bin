// 1bin / gorkscrew
// MIT License Copyright(c) 2020, 2022 Hiroshi Shimamoto
// vim:set sw=4 sts=4:

package gorkscrew

import (
    "io"
    "log"
    "os"

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
	log.Println(usage)
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
    if opt == "-F" && passfd == nil {
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
	passfd(conn)
	return
    }

    iorelay.Relay(conn, &stdio{ stdin: os.Stdin, stdout: os.Stdout })
}
