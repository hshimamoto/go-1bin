// 1bin / sshfwd
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package sshfwd

import (
    "io/ioutil"
    "log"
    "net"
    "strings"
    "time"

    "golang.org/x/crypto/ssh"
    "github.com/hshimamoto/go-iorelay"
    "github.com/hshimamoto/go-session"
)

func sshkeepalive(cli *ssh.Client) error {
    _, _, err := cli.SendRequest("keepalive@golang.org", true, nil)
    return err
}

func forwarder(cli *ssh.Client, lconn net.Conn, remote string) {
    defer lconn.Close()

    rconn, err := cli.Dial("tcp", remote)
    if err != nil {
	log.Printf("forwarder: Dial: %v\n", err)
	return
    }
    defer rconn.Close()

    log.Printf("start forwarding to %s\n", remote)
    iorelay.Relay(lconn, rconn)
    time.Sleep(time.Second)
    log.Printf("end forwarding to %s\n", remote)
}

type fwd struct {
    Local, Remote string
}

type req struct {
    remote string
    conn net.Conn
}

type host struct {
    Proxy, Dest, User, Key string
    Fwds []fwd
    //
    running bool
    keepalive time.Time
    q chan req
    nr int
}

func (h *host)Connect() (*ssh.Client, error) {
    cfg := &ssh.ClientConfig{
	User: h.User,
	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    if buf, err := ioutil.ReadFile(h.Key); err == nil {
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
	    log.Printf("ParsePrivateKey: %s %v\n", h.Key, err)
	    return nil, err
	}
	cfg.Auth = []ssh.AuthMethod{ ssh.PublicKeys(key) }
    } else {
	cfg.Auth = []ssh.AuthMethod{ ssh.Password(h.Key) }
    }

    // connect ssh
    var conn net.Conn
    var err error
    if h.Proxy == "" {
	conn, err = session.Dial(h.Dest)
	if err != nil {
	    log.Printf("Dial %s: %v\n", h.Dest, err)
	    return nil, err
	}
    } else {
	conn, err = session.Corkscrew(h.Proxy, h.Dest)
	if err != nil {
	    log.Printf("Corkscrew %s %s: %v\n", h.Proxy, h.Dest, err)
	    return nil, err
	}
    }
    // start ssh through conn
    cconn, cchans, creqs, err := ssh.NewClientConn(conn, h.Dest, cfg)
    if err != nil {
	log.Printf("NewClientConn %s: %v\n", h.Dest, err)
	conn.Close()
	return nil, err
    }
    cli := ssh.NewClient(cconn, cchans, creqs)
    log.Printf("ssh connection established with %s\n", h.Dest)
    return cli, nil
}

func (h *host)EventCheck(cli *ssh.Client) {
    select {
    case r := <-h.q:
	go func() {
	    h.nr++
	    forwarder(cli, r.conn, r.remote)
	    h.nr--
	}()
    case <-time.After(time.Minute):
	now := time.Now()
	if now.After(h.keepalive) {
	    if err := sshkeepalive(cli); err != nil {
		log.Printf("keepalive %v\n", err)
		h.running = false
		return
	    }
	    h.keepalive = now.Add(time.Minute)
	}
    }
}

func (h *host)RunForwarder() {
    for {
	if cli, err := h.Connect(); err == nil {
	    h.running = true
	    h.keepalive = time.Now().Add(time.Minute)
	    for h.running {
		h.EventCheck(cli)
	    }
	}
	// SSH session failed or done
	// interval 10s
	time.Sleep(10 * time.Second)
    }
}

func (h *host)RunListener(f fwd) {
    for {
	serv, err := session.NewServer(f.Local, func(conn net.Conn) {
	    if h.running {
		log.Printf("req forward %s to %s\n", f.Local, f.Remote)
		h.q <- req{ remote: f.Remote, conn: conn }
	    }
	})
	if err != nil {
	    log.Printf("NewServer %s %v\n", f.Local, err)
	    time.Sleep(time.Minute)
	    continue
	}
	serv.Run()
    }
}

func (h *host)Stats() {
    st := "stop"
    if h.running {
	st = "ruuning"
    }
    log.Printf("%s@%s: %s %d connections", h.User, h.Dest, st, h.nr)
}

func loadConfig(config string) []*host {
    hosts := []*host{}
    buf, err := ioutil.ReadFile(config)
    if err != nil {
	return []*host{}
    }
    var curr *host = nil
    proxy := ""
    for _, line := range strings.Split(string(buf), "\n") {
	if line == "" || line[0] == '#' {
	    continue
	}
	w := strings.Fields(line)
	if len(w) < 2 {
	    log.Printf("bad line: %s\n", line)
	    continue
	}
	switch w[0] {
	case "proxy": proxy = w[1]
	case "host":
	    if len(w) < 4 {
		log.Printf("bad line: %s\n", line)
		continue
	    }
	    curr = &host{
		Proxy: proxy,
		Dest: w[1],
		User: w[2],
		Key: w[3],
	    }
	    hosts = append(hosts, curr)
	case "fwd":
	    if len(w) < 3 {
		log.Printf("bad line: %s\n", line)
		continue
	    }
	    curr.Fwds = append(curr.Fwds, fwd{ Local: w[1], Remote: w[2] })
	}
    }

    return hosts
}

func Run(args []string) {
    if len(args) < 1 {
	log.Println("sshfwd <fwd config>")
	return
    }
    hosts := loadConfig(args[0])
    if len(hosts) == 0 {
	log.Println("no hosts")
	return
    }
    for _, host := range hosts {
	// initialize
	host.q = make(chan req)
	host.nr = 0
	// start sshfwd goroutines
	go host.RunForwarder()
	for _, fwd := range host.Fwds {
	    go host.RunListener(fwd)
	}
    }
    for {
	time.Sleep(time.Hour)
	for _, host := range hosts {
	    host.Stats()
	}
    }
}
