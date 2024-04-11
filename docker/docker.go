package docker

import (
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func StartDaemon() (out string, error string) {
	// TODO @suttont: Test this on Linux/Windows
	var cmdArr []string
	switch os := runtime.GOOS; os {
	case "darwin":
		cmdArr = []string{"open", "--background", "-a", "Docker"}
	case "windows":
		cmdArr = []string{"start-service", "docker"}
	}
	i := 0
	cmdFirst := cmdArr[i]
	cmdArr = append(cmdArr[:i], cmdArr[i+1:]...)
	cmd := exec.Command(cmdFirst, cmdArr...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err.Error()
	}
	return "Docker daemon started", ""
}

func StopDaemon() (out string, error string) {
	// TODO @suttont: Test on all platforms - generated code
	var cmdArr []string
	switch os := runtime.GOOS; os {
	case "darwin":
		cmdArr = []string{"osascript", "-e", `quit app "Docker"`}
	case "windows":
		cmdArr = []string{"stop-service", "docker"}
	}
	i := 0
	cmdFirst := cmdArr[i]
	cmdArr = append(cmdArr[:i], cmdArr[i+1:]...)
	cmd := exec.Command(cmdFirst, cmdArr...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err.Error()
	}
	return "Docker daemon stopped", ""
}

func Run() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", image.PullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
