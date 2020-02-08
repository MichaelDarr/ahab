package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// ContainerStatus returns a code corresponding to the status of a config's container. Codes:
// 0 - not found
// 1 - created
// 2 - restarting
// 3 - running
// 4 - removing
// 5 - paused
// 6 - exited
// 7 - dead
func ContainerStatus(config *Configuration, configPath string) (int, error) {
	output, err := DockerOutput(&[]string{"inspect", "-f", "'{{.State.Status}}'", ContainerName(config, configPath)})
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("Error checking container status: %s", err)
	}

	switch {
	case bytes.Contains(output, []byte("created")):
		return 1, nil
	case bytes.Contains(output, []byte("restarting")):
		return 2, nil
	case bytes.Contains(output, []byte("running")):
		return 3, nil
	case bytes.Contains(output, []byte("removing")):
		return 4, nil
	case bytes.Contains(output, []byte("paused")):
		return 5, nil
	case bytes.Contains(output, []byte("exited")):
		return 6, nil
	case bytes.Contains(output, []byte("dead")):
		return 7, nil
	default:
		return 0, fmt.Errorf("Unexpected container status: %s", err)
	}
}

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

// DockerContainerCmd runs a docker command on the active config's container
// opts is sequence of strings here because these commands are usually set statically in code
func DockerContainerCmd(config *Configuration, configPath string, cmd string, opts *[]string) error {
	containerOpts := []string{cmd, ContainerName(config, configPath)}
	if opts != nil {
		containerOpts = append(containerOpts, *opts...)
	}
	return DockerCmd(&containerOpts)
}

// DockerOutput runs a docker function behind the scenes and returns the output
func DockerOutput(opts *[]string) ([]byte, error) {
	cmd := Docker(opts)
	return cmd.Output()
}

// RemoveContainer removes an environment if it exists but is not running
func RemoveContainer(config *Configuration, configPath string) (err error) {
	status, err := ContainerStatus(config, configPath)
	PrintErrFatal(err)
	if status == 1 || status == 6 || status == 7 {
		err = DockerContainerCmd(config, configPath, "rm", nil)
	}
	return
}

// StopContainer stops an environment if it is running
func StopContainer(config *Configuration, configPath string) (err error) {
	status, err := ContainerStatus(config, configPath)
	PrintErrFatal(err)
	if status == 2 || status == 3 || status == 5 {
		err = DockerContainerCmd(config, configPath, "stop", nil)
	}
	return
}

// UpContainer creates and starts an environment
func UpContainer(config *Configuration, configPath string) (err error) {
	status, err := ContainerStatus(config, configPath)
	PrintErrFatal(err)
	if status == 0 {
		launchOpts := append([]string{"run", "-td"}, LaunchOpts(config, configPath)...)
		err = DockerCmd(&launchOpts)
	} else if status == 1 || status == 6 || status == 7 {
		err = DockerContainerCmd(config, configPath, "start", nil)
	} else if status == 5 {
		err = DockerContainerCmd(config, configPath, "unpause", nil)
	}
	return
}

// LaunchOpts returns options used to launch a container
func LaunchOpts(config *Configuration, configPath string) []string {
	opts := expandEnvs(&config.Options)
	for _, vol := range config.Volumes {
		opts = append(opts, "-v", prepVolumeString(vol, configPath))
	}

	if config.Workdir != "" {
		opts = append(opts, "-w", os.ExpandEnv(config.Workdir))
	}

	opts = append(opts, "--name", ContainerName(config, configPath))
	return append(opts, os.ExpandEnv(config.ImageURI))
}

// expandConfEnv expands environment variables present in a slice of strings
func expandEnvs(toExpand *[]string) []string {
	expanded := make([]string, len(*toExpand))
	for i, opt := range *toExpand {
		expanded[i] = os.ExpandEnv(opt)
	}
	return expanded
}

// prepVolumeString reformats a volume string, resolving local paths relative to the config dir
func prepVolumeString(rawVolume string, configPath string) string {
	// expand volume env vars and split by first ":" in string
	volumeSplit := strings.SplitN(os.ExpandEnv(rawVolume), ":", 2)

	// resolve first (local) path relative to config dir
	if strings.HasPrefix(volumeSplit[0], ".") {
		volumeSplit[0] = path.Join(filepath.Dir(configPath), strings.TrimLeft(volumeSplit[0], "."))
	} else if strings.HasPrefix(volumeSplit[0], "~") {
		volumeSplit[0] = path.Join(filepath.Dir(configPath), strings.TrimLeft(volumeSplit[0], "~"))
	} else if !strings.HasPrefix(volumeSplit[0], "/") {
		volumeSplit[0] = path.Join(filepath.Dir(configPath), volumeSplit[0], "~")
	}
	return strings.Join(volumeSplit, ":")
}
