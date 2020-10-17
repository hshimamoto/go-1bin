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
    "1bin/sshcompile"
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
    command{ name: "sshcompile", run: sshcompile.Run },
    command{ name: "sshfwd", run: sshfwd.Run },
    command{ name: "sshproxy", run: sshproxy.Run },
}

func header(cmd string, pid int) {
    log.Printf("1bin: multiple tools in one binary\n")
    log.Printf("1bin: <%s> pid=%d\n", cmd, pid)
}

func lookup(cmd string) *command {
    for _, e := range cmds {
	if e.name == cmd {
	    return &e
	}
    }
    return nil
}

func main() {
    cmd := filepath.Base(os.Args[0])
    args := os.Args[1:]
    pid := os.Getpid()

    // lookup cmd
    c := lookup(cmd)
    if c == nil {
	if len(args) > 0 {
	    cmd = args[0]
	    args = args[1:]
	    c = lookup(cmd)
	}
	if c == nil {
	    header(cmd, pid)
	    log.Printf("1bin: no command <%s>\n", cmd)
	    return
	}
    }

    // show header
    if !c.noheader {
	header(cmd, pid)
    }

    // setup logger
    log.SetFlags(log.Flags() | log.Lmsgprefix)
    log.SetPrefix(fmt.Sprintf("[%d <%s>] ", pid, cmd))

    // run command
    c.run(args)
}
