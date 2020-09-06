// 1bin / golangbuilder
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package golangbuilder

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/user"
    "path/filepath"
    "strings"

    "1bin/lib/f"
    "1bin/lib/docker"
)

func getgituser_from(file string) string {
    for line := range f.Lines(file) {
	w := strings.Fields(line)
	if len(w) < 3 {
	    continue
	}
	// w[0]  w[1] w[2]
	// email =    username@domain
	if w[0] == "email" {
	    return strings.Split(w[2], "@")[0]
	}
    }
    return ""
}

func getgituser() string {
    name := os.Getenv("USER")
    if u, err := user.Current(); err == nil {
	name = u.Username
    }
    home := os.Getenv("HOME")
    if n := getgituser_from(filepath.Join(home, ".gitconfig")); n != "" {
	name = n
    }
    if n := getgituser_from(".git/config"); n != "" {
	name = n
    }
    return name
}

func loadenvs_from(file string) ([]string, error) {
    envs := []string{}
    for line := range f.Lines(file) {
	if line == "" || line[0] == '#' {
	    continue
	}
	if strings.Index(line, "=") > 0 {
	    envs = append(envs, line)
	    continue
	}
	log.Printf("%s: syntax error\n", file)
	return envs, fmt.Errorf("syntax error");
    }
    return envs, nil
}

func loadenvs() []string {
    envfile := ".golangbuilder.env"
    envs := []string{}
    home := os.Getenv("HOME")
    if lines, err := loadenvs_from(filepath.Join(home, envfile)); err == nil {
	envs = append(envs, lines...)
    }
    if lines, err := loadenvs_from(envfile); err == nil {
	envs = append(envs, lines...)
    }
    return envs
}

func Run(args []string) {
    cwd, _ := os.Getwd()
    projname := filepath.Base(cwd)
    log.Printf("Dir: %s Project: %s\n", cwd, projname)
    gituser := getgituser()
    log.Printf("gituser: %s\n", gituser)
    envs := loadenvs()

    tempdir, err := ioutil.TempDir("", "golangbuilder")
    if err != nil {
	log.Printf("Tempdir: %v\n", err)
	return
    }
    defer os.RemoveAll(tempdir)
    log.Printf("TempDir: %s\n", tempdir)

    os.MkdirAll(filepath.Join(tempdir, "bin"), 0755)
    os.MkdirAll(filepath.Join(tempdir, "src"), 0755)

    // prepare HOME
    home := filepath.Join("src", "github.com", gituser)
    home_src := filepath.Join(tempdir, home)
    home_dst := filepath.Join("/go", home)
    os.MkdirAll(home_src, 0755)

    // working dir
    wd := filepath.Join(home_dst, projname)

    // setup docker
    d, err := docker.New("golang")
    if err != nil {
	log.Printf("docker %v\n", err)
	return
    }
    d.AddVol(tempdir, "/go")
    d.AddVol(cwd, wd)
    d.AddEnv("HOME", home_dst)
    d.WorkingDir(wd)
    for _, line := range envs {
	w := strings.SplitN(line, "=", 2)
	d.AddEnv(w[0], w[1])
    }

    // and run
    d.Run("bash")
}
