// 1bin / lib/f
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package f

import (
    "io/ioutil"
    "strings"
)

func Lines(file string) <-chan string {
    qlines := make(chan string)
    go func() {
	defer close(qlines)
	buf, err := ioutil.ReadFile(file)
	if err != nil {
	    return
	}
	for _, line := range strings.Split(string(buf), "\n") {
	    qlines <- line
	}
    }()
    return qlines
}
