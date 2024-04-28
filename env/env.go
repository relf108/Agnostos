package env

import (
	"context"
	"io"
	"os"
	"os/exec"

	"agnostos.com/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func CreateEnv(config config.Config) string {

	ctx := context.Background()

	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	imageStr := config.Lang.Name + ":" + config.Lang.Version

	reader, err := cli.ImagePull(ctx, imageStr, image.PullOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	// Create the container ready to be exec into - if the container stops immediately we can fix this by adding a `tail -f /dev/null` command
	containerConf := container.Config{
		Tty:   true,
		Image: imageStr,
		// TODO @suttont: Not connvinced we'll actually need this but it's here for now
		Cmd: []string{
			"tail -f /dev/null",
		},
		// TODO @suttont: This is currently failing, step 3 of this article may fix the issue https://gist.github.com/proudlygeek/5721498
		Env: []string{
			"SHARED_DIRECTORY=" + config.Mounts[0].Target,
			"FILEPERMISSIONS_UID=1000",
			"FILEPERMISSIONS_GID=1000",
			"FILEPERMISSIONS_MODE=777",
		},
	}

	hostConf := container.HostConfig{
		Mounts: _toMounts(config.Mounts),
	}

	resp, err := cli.ContainerCreate(ctx,
		&containerConf,
		&hostConf, nil, nil, "")
	if err != nil {
		panic(err)
	}

	return resp.ID
}

func EnterEnv(containerId string, cfg config.Config) {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(context.Background(), containerId, container.StartOptions{}); err != nil {
		panic(err)
	}

	cmd := exec.Command(
		"/bin/sh",
		[]string{
			"-c", "docker exec -it \"" + containerId + "\" \"/bin/sh\"",
		}...,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func MountInterpreters(containerId string, cfg config.Config) {
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(context.Background(), containerId, container.StartOptions{}); err != nil {
		panic(err)
	}

	//"docker exec -it \"" + containerId + "\" apt update -y && apt upgrade -y && apt-get install nfs-kernel-server portmap -y"
	cmd := exec.Command(
		"/bin/sh",
		[]string{
			"-c",
			"tail -f /dev/null",
		}...,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func _toMounts(mounts []*config.Mount) []mount.Mount {
	var res []mount.Mount
	for _, m := range mounts {
		res = append(res, mount.Mount{
			Type:   mount.TypeBind,
			Source: m.Source,
			Target: m.Target,
		})
	}
	return res
}
