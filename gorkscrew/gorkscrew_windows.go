// 1bin / gorkscrew
// MIT License Copyright(c) 2020, 2022 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package gorkscrew

import (
    "net"
)

var usage string = "gorkscrew <proxyhost:port> <remotehost:port>"
var passfd func(net.Conn) = nil
