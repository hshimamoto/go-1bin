// 1bin / bashhistcompact
// MIT License Copyright(c) 2020, 2021, 2022 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package bashhistcompact

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"
    "time"
)

func Run(args []string) {
    home, err := os.UserHomeDir()
    if err != nil {
	log.Printf("get HomeDir %v\n", err)
	return
    }

    histname := fmt.Sprintf("%s/.bash_history", home)
    // read everything
    hist, err := ioutil.ReadFile(histname)
    if err != nil {
	log.Printf("no .bash_history? %v\n", err)
	return
    }

    // try to create backup file ~/.bashhist
    pid := os.Getpid()
    // format YYYYMMDDhhmmss
    now := time.Now().Format("20060102150405")
    bkname := fmt.Sprintf("%s/.bashhist/histroy-%s-%d", home, now, pid)
    err = ioutil.WriteFile(bkname, hist, 0600)
    if err != nil {
	log.Printf("backup %s failed %v\n", bkname, err)
	return
    }

    // craft history
    f, err := os.OpenFile(histname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
    if err != nil {
	log.Fatal(err)
	return
    }
    defer f.Close()

    prev := ""
    for _, line := range strings.Split(string(hist), "\n") {
	switch line {
	case "ls", "pwd", "cd", "cd ..", "cd -": // simple command
	case "mount", "env", "top", "reset": // simple command
	case "ls -l", "ls -lh", "ls -lart", "ls -larth": // ls variant
	case "df", "df -h", "df -h .", "df -h ~": // df variant
	case "du -ksh .": // du variant
	case "ps", "ps auxww": // ps variant
	case "jobs", "%", "%1", "%2", "%3": // job control
	case "rm -rf *": // catastrophic
	case "make", "make clean": // simple make
	case prev: // duplicate
	case "": // empty
	default:
	    f.Write([]byte(line + "\n"))
	    prev = line
	}
    }
}
