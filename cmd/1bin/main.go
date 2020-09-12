// 1bin
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"

    "1bin/bashhistcompact"
    "1bin/fwd"
    "1bin/fwdset"
    "1bin/golangbuilder"
    "1bin/gradlebuilder"
    "1bin/sshfwd"
    "1bin/sshproxy"
)

type command struct {
    name string
    run func(args []string)
    noheader bool
}

var cmds = []command{
    command{ name: "bashhistcompact", run: bashhistcompact.Run, noheader: true },
    command{ name: "fwd", run: fwd.Run },
    command{ name: "fwdset", run: fwdset.Run },
    command{ name: "golangbuilder", run: golangbuilder.Run },
    command{ name: "gradlebuilder", run: gradlebuilder.Run },
    command{ name: "sshfwd", run: sshfwd.Run },
    command{ name: "sshproxy", run: sshproxy.Run },
}

func header(cmd string, pid int) {
    log.Printf("1bin: multiple tools in one binary\n")
    log.Printf("1bin: <%s> pid=%d\n", cmd, pid)
}

func main() {
    cmd := filepath.Base(os.Args[0])
    pid := os.Getpid()

    // lookup cmd
    var c *command = nil
    for _, e := range cmds {
	if e.name == cmd {
	    c = &e
	    break
	}
    }
    if c == nil {
	header(cmd, pid)
    }

    // show header
    if !c.noheader {
	header(cmd, pid)
    }

    // setup logger
    log.SetFlags(log.Flags() | log.Lmsgprefix)
    log.SetPrefix(fmt.Sprintf("[%d <%s>] ", pid, cmd))

    // run command
    c.run(os.Args[1:])
}
