package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ContainerUserName is the username of the user created in the container
const ContainerUserName = "ahab"

// Container contains all information regarding a container's configuration
type Container struct {
	Fields   *Configuration
	FilePath string
}

// GetContainer retrieves all container info relative to the working directory
func GetContainer() (*Container, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Failed to get working directory: %s", err)
	}

	container := new(Container)
	container.FilePath, err = findConfigPath(curDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to find config file: %s", err)
	}

	configFile, err := os.Open(container.FilePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open config file '%s': %s", container.FilePath, err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&container.Fields); err != nil {
		return nil, fmt.Errorf("Failed to parse config file '%s': %s", container.FilePath, err)
	}

	if missing := container.Fields.missingFields(); missing != "" {
		return nil, fmt.Errorf("Config file '%s' missing required fields: %s", container.FilePath, missing)
	}
	return container, checkConfigVersion(container.Fields.AhabVersion)
}

// Cmd runs a container command in the form docker [command] [container name]
func (container *Container) Cmd(command string) error {
	containerOpts := append([]string{command}, container.Name())
	return DockerCmd(&containerOpts)
}

// Create creates & prepares the container, leaving at "created" status if start is false
func (container *Container) Create(startContainer bool) error {
	launchOpts, err := container.creationOpts()
	if err != nil {
		return err
	}

	launchOpts = append([]string{"create", "-t"}, launchOpts...)
	if container.Fields.Command != "" {
		launchOpts = append(launchOpts, container.Fields.Command)
	} else {
		launchOpts = append(launchOpts, "top", "-b")
	}
	if err = DockerCmd(&launchOpts); err != nil {
		return err
	}

	// set up user permissions and start container
	if !container.Fields.Permissions.Disable {
		if err = container.Cmd("start"); err != nil {
			return err
		}
		if err = container.prep(); err != nil {
			return err
		}
	} else if startContainer || len(container.Fields.Init) > 0 {
		if err := container.Cmd("start"); err != nil {
			return err
		}
	}

	// execute init commands - if there are any, the container will have been started above
	for _, initCmd := range container.Fields.Init {
		initCmdSplit := rootExec(container, strings.Split(initCmd, " ")...)
		if err := DockerCmd(&initCmdSplit); err != nil {
			return nil
		}
	}

	// optionally stop/restart initial process after setup
	if s, _ := container.Status(); !startContainer && s == 3 {
		return container.Cmd("stop")
	} else if container.Fields.RestartAfterSetup {
		return container.Cmd("restart")
	}

	return nil
}

// Down stops and removes the container
func (container *Container) Down() (err error) {
	status, err := container.Status()
	if err == nil && (status == 2 || status == 3 || status == 5) {
		err = container.Cmd("stop")
		if err == nil {
			err = container.Cmd("rm")
		}
	} else if err == nil && (status == 1 || status == 6 || status == 7) {
		err = container.Cmd("rm")
	}
	return
}

// Name fetches the container name
func (container *Container) Name() string {
	if container.Fields.Name == "" {
		name := filepath.Dir(container.FilePath)
		name = strings.TrimPrefix(name, "/")
		return strings.ReplaceAll(name, "/", "_")
	}
	return container.Fields.Name
}

// Prop fetches a container property using Docker
func (container *Container) Prop(fieldID string) (string, error) {
	output, err := DockerOutput(&[]string{"inspect", "-f", "{{." + fieldID + "}}", container.Name()})
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return "", nil
	} else if err != nil {
		return "", err
	} else {
		idString := string(output)
		return strings.Trim(idString, " \n"), nil
	}
}

// Status returns a code corresponding to the status of the container
// 0 - not found
// 1 - created
// 2 - restarting
// 3 - running
// 4 - removing
// 5 - paused
// 6 - exited
// 7 - dead
func (container *Container) Status() (int, error) {
	status, err := container.Prop("State.Status")
	if err != nil {
		return 0, fmt.Errorf("Error checking container status: %s", err)
	}

	switch status {
	case "":
		return 0, nil
	case "created":
		return 1, nil
	case "restarting":
		return 2, nil
	case "running":
		return 3, nil
	case "removing":
		return 4, nil
	case "paused":
		return 5, nil
	case "exited":
		return 6, nil
	case "dead":
		return 7, nil
	default:
		return 0, fmt.Errorf("Unexpected container status: %s", status)
	}
}

// Up creates and starts the container
func (container *Container) Up() error {
	status, err := container.Status()
	if err != nil {
		return err
	} else if status == 0 {
		return container.Create(true)
	} else if status == 1 || status == 6 || status == 7 {
		return container.Cmd("start")
	} else if status == 5 {
		return container.Cmd("unpause")
	}
	return nil
}

// return a slice of options used when creating the container
func (container *Container) creationOpts() (opts []string, err error) {
	userConfig, err := UserConfig()
	if err != nil {
		return
	}

	// initial (idle) process runs as root user
	opts = []string{"-u", "root"}

	// user-specified options
	opts = append(opts, expandEnvs(&container.Fields.Options)...)
	opts = append(opts, expandEnvs(&userConfig.Options)...)

	// environment
	envStrings := expandEnvs(&container.Fields.Environment)
	envStrings = append(envStrings, expandEnvs(&userConfig.Environment)...)
	for _, envString := range envStrings {
		opts = append(opts, "-e", envString)
	}

	// volumes
	for _, vol := range container.Fields.Volumes {
		volString, err := prepVolumeString(vol, container.FilePath)
		if err != nil {
			return nil, err
		}
		opts = append(opts, "-v", volString)
	}
	for _, vol := range userConfig.Volumes {
		volString, err := prepVolumeString(vol, container.FilePath)
		if err != nil {
			return nil, err
		}
		opts = append(opts, "-v", volString)
	}

	// entrypoint
	if container.Fields.Entrypoint != "" {
		entrypointPath, err := prepVolumeString(container.Fields.Entrypoint, container.FilePath)
		if err != nil {
			return nil, err
		}
		opts = append(opts, []string{"--entrypoint", entrypointPath}...)
	}

	// workdir
	if container.Fields.Workdir != "" {
		opts = append(opts, "-w", os.ExpandEnv(container.Fields.Workdir))
	}

	// hostname
	if container.Fields.Hostname != "" {
		opts = append(opts, "-h", os.ExpandEnv(container.Fields.Hostname))
	} else {
		// hostname = name of parent dir of the container file
		opts = append(opts, "-h", filepath.Base(filepath.Dir(container.FilePath)))
	}

	// xhost sharing
	if container.Fields.ShareX11 {
		err = DockerXHostAuth()
		opts = append(opts, "-v", "/tmp/.X11-unix:/tmp/.X11-unix", "-e", "DISPLAY="+os.Getenv("DISPLAY"))
	}

	// container name and image
	return append(opts, "--name", container.Name(), os.ExpandEnv(container.Fields.ImageURI)), err
}

// run commands to prepare users/permissions in the container
func (container *Container) prep() error {
	homeDir := "/home/" + ContainerUserName
	uid := strconv.Itoa(os.Getuid())

	// commands which need to be executed after user creation
	// not using defer because it complicates error throwing
	extraCmds := [][]string{
		rootExec(container, "chown", ContainerUserName+":", homeDir),
		rootExec(container, "chmod", "700", homeDir),
	}

	// split out groups marked as "new" by a prefixed !, create them, and add the user (deferred)
	groups, newGroups := splitGroups(&container.Fields.Permissions.Groups)

	// create non-root user
	userAddCmd := rootExec(container)
	switch container.Fields.Permissions.CmdSet {
	case "", "default":
		userAddCmd = append(userAddCmd, []string{"useradd", "-o", "-m", "-d", homeDir, "-G"}...)
		userAddCmd = append(userAddCmd, strings.Join(groups, ","))
		for _, group := range newGroups {
			extraCmds = append(extraCmds, rootExec(container, "groupadd", group))
			extraCmds = append(extraCmds, rootExec(container, "usermod", "-G", group, ContainerUserName))
		}
	case "busybox":
		userAddCmd = append(userAddCmd, []string{"adduser", "-D", "-h", homeDir}...)
		for _, group := range newGroups {
			extraCmds = append(extraCmds, rootExec(container, "addgroup", group))
			extraCmds = append(extraCmds, rootExec(container, "addgroup", ContainerUserName, group))
		}
		for _, group := range groups {
			extraCmds = append(extraCmds, rootExec(container, "addgroup", ContainerUserName, group))
		}
	default:
		return fmt.Errorf("Unsupported command set specified in container: %s", container.Fields.Permissions.CmdSet)
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
