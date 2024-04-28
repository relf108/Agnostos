package main

import (
	// "agnostos.com/cli"
	"agnostos.com/config"
	"agnostos.com/docker"
	"agnostos.com/env"
)

func main() {
	docker.StartDaemon()
	// args := cli.ParseArgs()
	// println(
	// 	string(args.EnvOperator),
	// 	string(args.EnvName),
	// 	string(args.Lang.Name),
	// 	string(args.Lang.Version),
	// )
	path := config.FindConfig()
	config := config.ReadConfig(path)
	id := env.CreateEnv(config)
  // env.MountInterpreters(id, config)
	env.EnterEnv(id, config)
}
