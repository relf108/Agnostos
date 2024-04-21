package main

import (
	"agnostos.com/cli"
	"agnostos.com/docker"
	"agnostos.com/env"
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
  config := env.ReadConfig()
	id := env.CreateEnv(config)
	env.EnterEnv(id, config)
}
