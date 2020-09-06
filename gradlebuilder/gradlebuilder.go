// 1bin / gradlebuilder
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:
package gradlebuilder

import (
    "io/ioutil"
    "log"
    "os"
    "path/filepath"

    "1bin/lib/docker"
)

func Run(args []string) {
    cwd, _ := os.Getwd()
    projname := filepath.Base(cwd)
    log.Printf("Dir: %s Project: %s\n", cwd, projname)
    user := "gradle"

    tempdir, err := ioutil.TempDir("", "gradlebuilder")
    if err != nil {
	log.Printf("Tempdir: %v\n", err)
	return
    }
    defer os.RemoveAll(tempdir)
    log.Printf("TempDir: %s\n", tempdir)

    // prepare HOME
    os.MkdirAll(filepath.Join(tempdir, "home", user, projname), 0755)

    // working dir
    wd := filepath.Join("/home", user, projname)

    // setup docker
    d, err := docker.New("gradle", "jdk14")
    if err != nil {
	log.Printf("docker %v\n", err)
	return
    }
    d.AddVol(filepath.Join(tempdir, "home"), "/home")
    d.AddVol(cwd, wd)
    d.AddEnv("HOME", filepath.Join("/home", user))
    d.WorkingDir(wd)

    // and run
    d.Run("bash")
}
