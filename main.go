package main

import (
	"agnostos.com/cli"
	"agnostos.com/docker"
)

func main() {
	out, err := docker.StartDaemon()
	if err != "" {
		panic(err)
	}
	println(out)
	cli.ReadArgs()
	docker.Run()
}
