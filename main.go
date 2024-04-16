package main

import (
	"agnostos.com/cli"
	"agnostos.com/docker"
)

func main() {
	docker.StartDaemon()
	args := cli.ParseArgs()
	println(
		string(args.EnvOperator),
		string(args.EnvName),
		string(args.Lang.Name),
		string(args.Lang.Version),
	)
docker.Run()
}
