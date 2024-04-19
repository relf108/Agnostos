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
)

// TODO @suttont: Test this function on all OSs
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

// TODO @suttont: Test this function on all OSs
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
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Tty:   true,
		Image: imageStr,
	}, nil, nil, nil, "")
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
