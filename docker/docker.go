package docker

import (
	"os"
	"os/exec"
	"runtime"
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
