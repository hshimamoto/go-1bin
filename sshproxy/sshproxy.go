// 1bin / sshproxy
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package sshproxy

import (
    "bytes"
    "log"
    "net"
    "time"
    "strings"

    "github.com/hshimamoto/go-session"
    "github.com/hshimamoto/go-iorelay"
    "github.com/BurntSushi/toml"
)

// Example config.toml
//
// Listen = "127.0.0.1:8822"
// Upstream = ""
//
// [[Hosts]]
// Hostname = "host1"
// Dest = "host1:22"
//
// [[Hosts]]
// Hostname = "host2"
// Dest = "host2:22"

type host struct {
    Hostname, Dest, Proxy string
}

type proxyconfig struct {
    Listen, Upstream string
    Hosts []host
}

func loadconfig(config string) (*proxyconfig, error) {
    cfg := &proxyconfig{}
    if _, err := toml.DecodeFile(config, cfg); err != nil {
	return nil, err
    }
    return cfg, nil
}

func handle(cfg *proxyconfig, conn net.Conn) {
    buf := make([]byte, 256)
    rest := []byte(nil)
    n := 0
    // get entire CONNECT header
    for {
	r, err := conn.Read(buf[n:])
	if err != nil {
	    log.Printf("read header error: %v\n", err)
	    return
	}
	if r == 0 {
	    log.Println("connection closed")
	    return
	}
	n += r
	if idx := bytes.Index(buf[0:n], []byte{13, 10, 13, 10}); idx > 0 {
	    if n > idx + 4 {
		rest = buf[idx + 4:n]
	    }
	    break
	}
	if n >= 256 {
	    log.Println("read header too large")
	    return
	}
    }
    // try to connect upstream
    lines := strings.Split(string(buf[0:n]), "\r\n")
    // CONNECT host:port HTTP/1.0
    w := strings.Fields(lines[0])
    if len(w) < 3 {
	log.Printf("request error: %s\n", lines[0])
	return
    }
    if w[0] != "CONNECT" {
	log.Printf("no connect request: %s\n", lines[0])
	return
    }
    hostport := strings.Split(w[1], ":")
    if len(hostport) != 2 {
	log.Printf("invalid hostport: %s\n", lines[0])
	return
    }
    host := hostport[0]
    port := hostport[1]
    log.Printf("CONNECTing to %s:%s\n", host, port)
    dest := ""
    dest2 := ""
    // lookup host
    for _, h := range cfg.Hosts {
	if h.Hostname == host {
	    if h.Proxy != "" {
		dest = h.Proxy
		dest2 = h.Dest
	    } else {
		dest = h.Dest
	    }
	    break
	}
    }
    if dest == "" {
	log.Printf("no host %s found\n", host)
	return
    }
    var fconn net.Conn = nil
    var err error
    if cfg.Upstream != "" {
	fconn, err = session.Corkscrew(cfg.Upstream, dest)
	if err != nil {
	    log.Printf("Corkscrew: %s %s %v\n", cfg.Upstream, dest, err)
	    return
	}
    } else {
	fconn, err = session.Dial(dest)
	if err != nil {
	    log.Printf("Dial: %s %v\n", dest, err)
	    return
	}
    }
    // connection established
    defer fconn.Close()
    if dest2 != "" {
	log.Printf("HttpConnect: %s %s\n", dest, dest2)
	if err := session.HttpConnect(fconn, dest2); err != nil {
	    log.Printf("HttpConnect: %v\n", err)
	    return
	}
    }
    log.Printf("connection %s:%s established\n", host, port)
    conn.Write([]byte("HTTP/1.1 200 Established\r\n\r\n"))
    if rest != nil {
	log.Printf("rest %s\n", string(rest))
	fconn.Write(rest)
    }
    iorelay.RelayWithTimeout(conn, fconn, time.Hour)
}

func Run(args []string) {
    if len(args) < 1 {
	log.Println("sshproxy <config toml>")
	return
    }
    config := args[0]
    cfg, err := loadconfig(config)
    if err != nil {
	log.Printf("config %s error: %v", config, err)
	return
    }
    // now we have Listen address
    serv, err := session.NewServer(cfg.Listen, func(conn net.Conn) {
	defer conn.Close()
	cfg, err := loadconfig(config)
	if err != nil {
	    log.Printf("config %s error: %v", config, err)
	    return
	}
	handle(cfg, conn)
    })
    if err != nil {
	log.Printf("NewServer error: %v\n", err)
	return
    }
    serv.Run()
}
