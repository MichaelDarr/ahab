package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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

// versionOrdinal standardizes version strings for easy comparison
// see https://stackoverflow.com/a/18411978
func versionOrdinal(version string) (string, error) {
	// ISO/IEC 14651:2011
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			return "", fmt.Errorf("versionOrdinal: Invalid version '%s'", version)
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo), nil
}
