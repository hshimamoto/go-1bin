// 1bin / fwdset
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package fwdset

import (
    "log"
    "os"
    "os/exec"
    "time"

    "github.com/BurntSushi/toml"
)

// Example config.toml
//
// [[Fwds]]
// Name = "name"
// Src = ":8080"
// Dst = "proxy:8080"

type fwd struct {
    Name, Src, Dst string
}

type fwdconfig struct {
    Fwds []fwd
}

func loadconfig(config string) (*fwdconfig, error) {
    cfg := &fwdconfig{}
    if _, err := toml.DecodeFile(config, cfg); err != nil {
	log.Printf("config %s error: %v", config, err)
	return nil, err
    }
    return cfg, nil
}

type fwdproc struct {
    fwd
    cmd *exec.Cmd
}

func (fp *fwdproc)Start(q chan *fwdproc) {
    log.Printf("start fwd %s %s %s\n", fp.Name, fp.Src, fp.Dst)
    fp.cmd = exec.Command("fwd", fp.Src, fp.Dst)
    fp.cmd.Stdout = os.Stdout
    fp.cmd.Stderr = os.Stderr
    go func() {
	err := fp.cmd.Run()
	log.Printf("done fwd %s %s %s %v\n", fp.Name, fp.Src, fp.Dst, err)
	time.Sleep(time.Second) // interval
	q <- fp // send signal
	fp.cmd = nil
    }()
}

func manage(cfg *fwdconfig) {
    // setup fwdprocs
    fps := []*fwdproc{}
    for _, fwd := range cfg.Fwds {
	fps = append(fps, &fwdproc{ fwd: fwd, cmd: nil })
    }
    q := make(chan *fwdproc)
    for {
	for _, fp := range fps {
	    if fp.cmd == nil {
		fp.Start(q)
	    }
	}
	// wait a proc done
	<-q
    }
}

func Run(args []string) {
    if len(args) < 1 {
	log.Println("fwdset <config toml>")
	return
    }
    cfg, err := loadconfig(args[0])
    if err != nil {
	return
    }
    manage(cfg)
}
