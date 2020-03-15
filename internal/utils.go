package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// DisplaySessionType returns the host's session type (x11, wayland, etc.)
func DisplaySessionType() (sessionType string) {
	// The env var XDG_SESSION_TYPE should tell us this, but some wayland compositors (sway) don't
	// set it on launch and it remains "tty". So, if it's "tty", we check if WAYLAND_DISPLAY is set
	// to detect wayland.
	sessionType = os.Getenv("XDG_SESSION_TYPE")
	_, waylandPresent := os.LookupEnv("WAYLAND_DISPLAY")
	if sessionType == "tty" && waylandPresent {
		sessionType = "wayland"
	}
	return
}

// ParseStatus returns the readable string of the passed container status
// 0 - not found
// 1 - created
// 2 - restarting
// 3 - running
// 4 - removing
// 5 - paused
// 6 - exited
// 7 - dead
func ParseStatus(code int) string {
	switch code {
	case 0:
		return "not found"
	case 1:
		return "created"
	case 2:
		return "restarting"
	case 3:
		return "running"
	case 4:
		return "removing"
	case 5:
		return "paused"
	case 6:
		return "exited"
	case 7:
		return "dead"
	default:
		return "unknown"
	}
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

func waylandDisplayOptions() []string {
	return []string{
		"-v", os.ExpandEnv("$XDG_RUNTIME_DIR/$WAYLAND_DISPLAY:/tmp/host-wayland"),
		"-e", "XDG_RUNTIME_DIR=/tmp", "-e", "WAYLAND_DISPLAY=host-wayland",
		"-e", "CLUTTER_BACKEND=wayland", "-e", "GDK_BACKEND=wayland",
		"-e", "QT_QPA_PLATFORM=wayland", "-e", "DL_VIDEODRIVER=wayland",
		"-e", "ELM_DISPLAY=wl", "-e", "ELM_ACCEL=opengl", "-e", "ECORE_EVAS_ENGINE=wayland_egl"}
}

func xDisplayOptions() []string {
	return []string{
		"-v", "/tmp/.X11-unix:/tmp/.X11-unix",
		"-e", "DISPLAY=" + os.Getenv("DISPLAY")}
}
