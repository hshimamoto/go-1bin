// 1bin / lib/docker
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package docker

import (
    "fmt"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
)

type Docker struct {
    image string
    cname string
    ugid string
    wd string
    vols []string
    envs []string
}

func New(name, tag string) (*Docker, error) {
    // container name
    cwd, _ := os.Getwd()
    dirname := filepath.Base(cwd)
    cname := fmt.Sprintf("%s-%s-%d", name, dirname, os.Getpid())

    // get user
    u, err := user.Current()
    if err != nil {
	return nil, err
    }
    ugid := u.Uid + ":" + u.Gid

    d := &Docker{}
    d.image = name + ":" + tag
    d.cname = cname
    d.ugid = ugid
    d.vols = []string{}
    d.envs = []string{}
    return d, nil
}

func (d *Docker)AddVol(src, dst string) {
    d.vols = append(d.vols, fmt.Sprintf("%s:%s", src, dst))
}

func (d *Docker)AddEnv(name, val string) {
    d.envs = append(d.envs, fmt.Sprintf("%s=%s", name, val))
}

func (d *Docker)WorkingDir(dir string) {
    d.wd = dir
}

func (d *Docker)Run(cmdline string) error {
    cmd := exec.Command(
	"docker", "run", "-it", "--rm",
	"--name", d.cname, "--hostname", d.cname,
	"-u", d.ugid)
    for _, v := range d.vols {
	cmd.Args = append(cmd.Args, "-v", v)
    }
    for _, e := range d.envs {
	cmd.Args = append(cmd.Args, "-e", e)
    }
    if d.wd != "" {
	cmd.Args = append(cmd.Args, "-w", d.wd)
    }
    cmd.Args = append(cmd.Args, d.image, cmdline)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return err
}
