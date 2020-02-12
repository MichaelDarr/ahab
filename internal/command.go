package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// ContainerUserName is the username of the user created in the container
const ContainerUserName = "ahab"

// ContainerProp returns a property of a config's container, or an empty string if it is not created
func ContainerProp(config *Configuration, configPath string, fieldID string) (string, error) {
	output, err := DockerOutput(&[]string{"inspect", "-f", "{{." + fieldID + "}}", ContainerName(config, configPath)})
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return "", nil
	} else if err != nil {
		return "", err
	} else {
		idString := string(output)
		return strings.Trim(idString, " \n"), nil
	}
}

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
	output, err := DockerOutput(&[]string{"inspect", "-f", "{{.State.Status}}", ContainerName(config, configPath)})
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
func DockerContainerCmd(config *Configuration, configPath string, cmd string) error {
	containerOpts := append([]string{cmd}, ContainerName(config, configPath))
	return DockerCmd(&containerOpts)
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

// LaunchOpts prepares the host to launch a container and returns options used to launch it
func LaunchOpts(config *Configuration, configPath string) (opts []string, err error) {
	userConfig, err := UserConfig()

	// initial (idle) process runs as root user
	opts = []string{"-u", "root"}

	// user-specified options
	opts = append(opts, expandEnvs(&config.Options)...)
	opts = append(opts, expandEnvs(&userConfig.Options)...)

	// environment
	envStrings := expandEnvs(&config.Environment)
	envStrings = append(envStrings, expandEnvs(&userConfig.Environment)...)
	for _, envString := range envStrings {
		opts = append(opts, "-e", envString)
	}

	// volumes
	for _, vol := range config.Volumes {
		volString, err := prepVolumeString(vol, configPath)
		if err != nil {
			return nil, err
		}
		opts = append(opts, "-v", volString)
	}
	for _, vol := range userConfig.Volumes {
		volString, err := prepVolumeString(vol, configPath)
		if err != nil {
			return nil, err
		}
		opts = append(opts, "-v", volString)
	}

	// entrypoint
	if config.Entrypoint != "" {
		entrypointPath, err := prepVolumeString(config.Entrypoint, configPath)
		if err != nil {
			return nil, err
		}
		opts = append(opts, []string{"--entrypoint", entrypointPath}...)
	}

	// workdir
	if config.Workdir != "" {
		opts = append(opts, "-w", os.ExpandEnv(config.Workdir))
	}

	// hostname
	if config.Hostname != "" {
		opts = append(opts, "-h", config.Hostname)
	} else {
		opts = append(opts, "-h", filepath.Base(filepath.Dir(configPath))) // name of parent dir of the config file
	}

	// xhost sharing
	if config.ShareX11 {
		err = DockerXHostAuth()
		opts = append(opts, "-v", "/tmp/.X11-unix:/tmp/.X11-unix", "-e", "DISPLAY="+os.Getenv("DISPLAY"))
	}

	// container name and image
	return append(opts, "--name", ContainerName(config, configPath), os.ExpandEnv(config.ImageURI)), err
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

// CreateContainer creates the container described by a config
func CreateContainer(config *Configuration, configPath string, startContainer bool) error {
	launchOpts, err := LaunchOpts(config, configPath)
	if err != nil {
		return err
	}

	launchOpts = append([]string{"create", "-t"}, launchOpts...)
	if config.Command != "" {
		launchOpts = append(launchOpts, config.Command)
	} else {
		launchOpts = append(launchOpts, "top", "-b")
	}
	if err = DockerCmd(&launchOpts); err != nil {
		return err
	}

	// set up user permissions and start container
	if !config.Permissions.Disable {
		if err = DockerContainerCmd(config, configPath, "start"); err != nil {
			return err
		}
		if err = PrepContainer(config, configPath); err != nil {
			return err
		}
	} else if startContainer || len(config.Init) > 0 {
		if err := DockerContainerCmd(config, configPath, "start"); err != nil {
			return err
		}
	}

	// execute init commands - if there are any, the container will have been started above
	for _, initCmd := range config.Init {
		initCmdSplit := rootExec(config, configPath, strings.Split(initCmd, " ")...)
		if err := DockerCmd(&initCmdSplit); err != nil {
			return nil
		}
	}

	// optionally stop/restart initial process after setup
	if s, _ := ContainerStatus(config, configPath); !startContainer && s == 3 {
		return DockerContainerCmd(config, configPath, "stop")
	} else if config.RestartAfterSetup {
		return DockerContainerCmd(config, configPath, "restart")
	}

	return nil
}

// PrepContainer executes docker functions to prepare a non-root user in the container
func PrepContainer(config *Configuration, configPath string) error {
	homeDir := "/home/" + ContainerUserName
	uid := strconv.Itoa(os.Getuid())

	// commands which need to be executed after user creation
	// not using defer because it complicates error throwing
	extraCmds := [][]string{
		rootExec(config, configPath, "chown", ContainerUserName+":", homeDir),
		rootExec(config, configPath, "chmod", "700", homeDir),
	}

	// split out groups marked as "new" by a prefixed !, create them, and add the user (deferred)
	groups, newGroups := splitGroups(&config.Permissions.Groups)

	// create non-root user
	userAddCmd := rootExec(config, configPath)
	switch config.Permissions.CmdSet {
	case "", "default":
		userAddCmd = append(userAddCmd, []string{"useradd", "-o", "-d", homeDir, "-G"}...)
		userAddCmd = append(userAddCmd, strings.Join(groups, ","))
		for _, group := range newGroups {
			extraCmds = append(extraCmds, rootExec(config, configPath, "groupadd", group))
			extraCmds = append(extraCmds, rootExec(config, configPath, "usermod", "-G", group, ContainerUserName))
		}
	case "busybox":
		userAddCmd = append(userAddCmd, []string{"adduser", "-D", "-h", homeDir}...)
		for _, group := range newGroups {
			extraCmds = append(extraCmds, rootExec(config, configPath, "addgroup", group))
			extraCmds = append(extraCmds, rootExec(config, configPath, "addgroup", ContainerUserName, group))
		}
		for _, group := range groups {
			extraCmds = append(extraCmds, rootExec(config, configPath, "addgroup", ContainerUserName, group))
		}
	default:
		return fmt.Errorf("Unsupported command set specified in config: %s", config.Permissions.CmdSet)
	}
	userAddCmd = append(userAddCmd, []string{"-u", uid, ContainerUserName}...)
	if err := DockerCmd(&userAddCmd); err != nil {
		return err
	}

	// run post-user-creation commands
	for _, userCmd := range extraCmds {
		if err := DockerCmd(&userCmd); err != nil {
			return err
		}
	}
	return nil
}

// RemoveContainer removes an environment if it exists but is not running
func RemoveContainer(config *Configuration, configPath string) error {
	status, err := ContainerStatus(config, configPath)
	if err == nil && (status == 1 || status == 6 || status == 7) {
		return DockerContainerCmd(config, configPath, "rm")
	}
	return err
}

// StopContainer stops an environment if it is running
func StopContainer(config *Configuration, configPath string) error {
	status, err := ContainerStatus(config, configPath)
	if err == nil && (status == 2 || status == 3 || status == 5) {
		return DockerContainerCmd(config, configPath, "stop")
	}
	return err
}

// UpContainer creates and starts an environment
func UpContainer(config *Configuration, configPath string) error {
	status, err := ContainerStatus(config, configPath)
	if err != nil {
		return err
	} else if status == 0 {
		return CreateContainer(config, configPath, true)
	} else if status == 1 || status == 6 || status == 7 {
		return DockerContainerCmd(config, configPath, "start")
	} else if status == 5 {
		return DockerContainerCmd(config, configPath, "unpause")
	}
	return nil
}

// expandConfEnv expands environment variables present in a slice of strings
func expandEnvs(toExpand *[]string) []string {
	expanded := make([]string, len(*toExpand))
	for i, opt := range *toExpand {
		expanded[i] = os.ExpandEnv(opt)
	}
	return expanded
}

// prepVolumeString resolves the first local path in a string (before ":") relative to the config dir
func prepVolumeString(rawVolume string, configPath string) (string, error) {
	// expand volume env vars and split by first ":" in string
	volumeSplit := strings.SplitN(os.ExpandEnv(rawVolume), ":", 2)

	// resolve first (local) path relative to config dir
	if strings.HasPrefix(volumeSplit[0], "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		volumeSplit[0] = path.Join(homeDir, strings.TrimLeft(volumeSplit[0], "~"))
	} else if !strings.HasPrefix(volumeSplit[0], "/") {
		volumeSplit[0] = path.Join(filepath.Dir(configPath), volumeSplit[0])
	}
	return strings.Join(volumeSplit, ":"), nil
}

// rootExec transforms a command into a DockerCmd-ready root-executed command
func rootExec(config *Configuration, configPath string, args ...string) []string {
	return append([]string{"exec", "-u", "root", ContainerName(config, configPath)}, args...)
}

// some groups are prefixed with ! - these are "new groups". splitGroups divides these out and removes the prefix
func splitGroups(allGroups *[]string) (groups []string, newGroups []string) {
	for _, group := range *allGroups {
		if strings.HasPrefix(group, "!") {
			newGroups = append(newGroups, strings.TrimLeft(group, "!"))
		} else {
			groups = append(groups, group)
		}
	}
	return
}
