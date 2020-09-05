// 1bin
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"

    "1bin/fwd"
    "1bin/fwdset"
)

func main() {
    cmd := filepath.Base(os.Args[0])
    pid := os.Getpid()
    log.Printf("1bin: %d <%s>\n", pid, cmd)
    log.SetFlags(log.Flags() | log.Lmsgprefix)
    log.SetPrefix(fmt.Sprintf("[%d <%s>] ", pid, cmd))
    switch cmd {
    case "fwd": fwd.Run(os.Args[1:])
    case "fwdset": fwdset.Run(os.Args[1:])
    }
}
