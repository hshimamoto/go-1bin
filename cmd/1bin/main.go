// 1bin
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package main

import (
    "log"
    "os"
    "path/filepath"

    "1bin/fwd"
)

func main() {
    cmd := filepath.Base(os.Args[0])
    log.Printf("1bin: %s\n", cmd)
    switch cmd {
    case "fwd": fwd.Run(os.Args[1:])
    }
}
