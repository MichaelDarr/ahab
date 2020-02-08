package internal

// LaunchOpts returns a slice of options used to launch a container
func LaunchOpts(config *Configuration) (opts []string, err error) {
	opts = append(config.Options, config.ImageURI)
	return
}
