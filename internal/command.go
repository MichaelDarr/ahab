package internal

import (
	"os"
	"os/exec"
)

// Docker generates a Docker command
func Docker(opts *[]string) (cmd *exec.Cmd) {
	cmd = exec.Command("docker", *opts...)
	return
}

// DockerCmd prints a function, runs it, and attatches the output to the user's terminal
func DockerCmd(opts *[]string) error {
	cmd := Docker(opts)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	PrintCmd(cmd)
	return cmd.Run()
}

// DockerOutput runs a docker function behind the scenes and returns the output
func DockerOutput(opts *[]string) ([]byte, error) {
	cmd := Docker(opts)
	return cmd.Output()
}

// DockerXHostAuth gives the docker user xhost access so containerized apps can run on the host's WM
// TODO: check if the Docker user already has access before re-granting access
func DockerXHostAuth() error {
	cmd := exec.Command("xhost", "+local:docker")
	return cmd.Run()
}

// ListContainers executes a docker command to list all containers
func ListContainers(verbose bool) error {
	execArgs := []string{"ps", "-a"}
	if !verbose {
		execArgs = append(execArgs, "--format", "table {{.Names}}\t{{.ID}}\t{{.Status}}")
	}
	return DockerCmd(&execArgs)
}

// ListImages executes a docker command to list all images
func ListImages(verbose bool) error {
	execArgs := []string{"images"}
	if !verbose {
		execArgs = append(execArgs, "--format", "table {{.Tag}}\t{{.ID}}\t{{.Size}}")
	}
	return DockerCmd(&execArgs)
}

// ListVolumes executes a docker command to list all images
func ListVolumes() error {
	execArgs := []string{"volume", "ls"}
	return DockerCmd(&execArgs)
}

// rootExec transforms a command into a DockerCmd-ready root-executed command
func rootExec(config *Configuration, configPath string, args ...string) []string {
	return append([]string{"exec", "-u", "root", ContainerName(config, configPath)}, args...)
}
