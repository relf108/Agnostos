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

// TODO @suttont: This should be for passing a config file, args are here as a placeholder

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

	cmd := exec.Command("/bin/sh", []string{"-c", "docker exec -it \"" + containerId + "\" \"/bin/sh\""}...)
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
