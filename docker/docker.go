package docker

import (
	"context"
	"os"
	"os/exec"
	"runtime"

	// "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func StartDaemon() {
	var cmdArr []string

	switch os := runtime.GOOS; os {
	case "darwin":
		cmdArr = []string{"open", "--background", "-a", "Docker"}
	case "windows":
		cmdArr = []string{"start-service", "docker"}
	case "linux":
		cmdArr = []string{"systemctl", "start", "docker"}
	}

	i := 0
	cmdFirst := cmdArr[i]
	cmdArr = append(cmdArr[:i], cmdArr[i+1:]...)
	cmd := exec.Command(cmdFirst, cmdArr...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		panic(err)
	}
	println("Docker daemon started")
}

func StopDaemon() (out string, error string) {
	var cmdArr []string

	switch os := runtime.GOOS; os {
	case "darwin":
		cmdArr = []string{"osascript", "-e", `quit app "Docker"`}
	case "windows":
		cmdArr = []string{"stop-service", "docker"}
	case "linux":
		cmdArr = []string{"systemctl", "stop", "docker"}
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

	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// Pull the latest Ubuntu image
	_, err = cli.ImagePull(ctx, "docker.io/library/ubuntu:latest", image.PullOptions{})

	// Create the container ready to be exec into - if the container stops immediately we can fix this by adding a `tail -f /dev/null` command
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Tty:   true,
		Image: "python:latest",
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	// use normal cli commands to exec into the container
	cmd := exec.Command("/bin/sh", []string{"-c", "docker exec -it \"" + resp.ID + "\" \"/bin/sh\""}...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}
