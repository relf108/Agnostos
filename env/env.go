package env

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/apple/pkl-go/pkl"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// TODO @suttont: This should be for passing a config file, args are here as a placeholder
type NosConfig struct {
	EnvName     string `pkl:"env_name"`
	LangName    string `pkl:"lang_name"`
	LangVersion string `pkl:"lang_version"`
}

func CreateEnv(envName string, langName string, langVersion string) string {

	ctx := context.Background()

	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	imageStr := langName + ":" + langVersion

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
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/Users/tristan.sutton/Projects/Agnostos", // TODO @suttont: This should be read from config
				Target: "/path/in/container",                      // TODO @suttont: This should be read from config
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx,
		&containerConf,
		&hostConf, nil, nil, "")
	if err != nil {
		panic(err)
	}

	return resp.ID
}

func EnterEnv(containerId string) {
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

func ReadConfig() {
	evaluator, err := pkl.NewEvaluator(context.Background(), pkl.PreconfiguredOptions)
	if err != nil {
		panic(err)
	}
	defer evaluator.Close()
	var cfg NosConfig
	if err = evaluator.EvaluateModule(context.Background(), pkl.FileSource("foo.pkl"), &cfg); err != nil {
		panic(err)
	}
}
