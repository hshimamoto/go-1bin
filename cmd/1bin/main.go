// 1bin
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package main

import (
    "log"
    "os"
    "path/filepath"
)

func main() {
    cmd := filepath.Base(os.Args[0])
    log.Println(cmd)
}
