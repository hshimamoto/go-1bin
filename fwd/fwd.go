// 1bin / fwd
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package fwd

import (
    "log"
    "net"
    "strconv"
    "time"

    "github.com/hshimamoto/go-session"
    "github.com/hshimamoto/go-iorelay"
    "github.com/mxk/go-flowrate/flowrate"
)

type FlowrateIO struct {
    r *flowrate.Reader
    w *flowrate.Writer
}

func (f *FlowrateIO)Read(p []byte) (int, error) {
    return f.r.Read(p)
}

func (f *FlowrateIO)Write(p []byte) (int, error) {
    return f.w.Write(p)
}

func Run(args []string) {
    if len(args) < 2 {
	log.Println("fwd <listen> <dst> [limit KB/s or MB/s]")
	return
    }
    limit := int64(0)
    if len(args) > 2 {
	lim := args[2]
	unit := 0
	switch lim[len(lim) - 1] {
	case 'k', 'K': unit = 1024
	case 'm', 'M': unit = 1024 * 1024
	}
	num, err := strconv.Atoi(lim[:len(lim) - 1])
	if err != nil {
	    num = 0
	}
	limit = int64(num * unit)
	if limit == 0 {
	    log.Printf("bad option: %s\n", lim)
	    return
	}
	log.Printf("set Ratelimit %d bytes/s", limit)
    }
    timeout := time.Hour
    serv, err := session.NewServer(args[0], func(conn net.Conn) {
	defer conn.Close()
	fconn, err := session.Dial(args[1])
	if err != nil {
	    log.Printf("Dial %s %v\n", args[1], err)
	    return
	}
	defer fconn.Close()
	if limit > 0 {
	    f1 := &FlowrateIO{
		r: flowrate.NewReader(conn, limit),
		w: flowrate.NewWriter(conn, limit),
	    }
	    f2 := &FlowrateIO{
		r: flowrate.NewReader(fconn, limit),
		w: flowrate.NewWriter(fconn, limit),
	    }
	    iorelay.RelayWithTimeout(f1, f2, timeout)
	} else {
	    iorelay.RelayWithTimeout(conn, fconn, timeout)
	}
    })
    if err != nil {
	log.Printf("NewServer %s %v\n", args[0], err)
	return
    }
    serv.Run()
}
