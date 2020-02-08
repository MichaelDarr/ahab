package internal

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Docker prints and runs a Docker command
func Docker(opts *[]string) (err error) {
	cmd := exec.Command("docker", *opts...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	PrintCmd(cmd)
	return cmd.Run()
}

// LaunchOpts returns options used to launch a container
func LaunchOpts(config *Configuration, configPath string) []string {
	opts := expandEnvs(&config.Options)
	for _, vol := range config.Volumes {
		opts = append(opts, "-v", prepVolumeString(vol, configPath))
	}

	for _, env := range expandEnvs(&config.Environment) {
		opts = append(opts, "-e", env)
	}

	if config.Workdir != "" {
		opts = append(opts, "-w", os.ExpandEnv(config.Workdir))
	}

	if config.Name != "" {
		opts = append(opts, "--name", config.Name)
	} else {
		opts = append(opts, "--name", ContainerPathName(configPath))
	}

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
