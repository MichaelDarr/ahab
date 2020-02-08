package internal

import "os"

// LaunchOpts returns options used to launch a container
func LaunchOpts(config *Configuration) (opts []string, err error) {
	opts = expandEnvs(&config.Options)
	opts = append(opts, os.ExpandEnv(config.ImageURI))
	return
}

// expandConfEnv expands environment variables present in a slice of strings
func expandEnvs(toExpand *[]string) []string {
	expanded := make([]string, len(*toExpand))
	for i, opt := range *toExpand {
		expanded[i] = os.ExpandEnv(opt)
	}
	return expanded
}
