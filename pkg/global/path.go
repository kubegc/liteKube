package global

import (
	"os/user"
	"path/filepath"
)

var HomePath string = GetHomeDir()

func GetHomeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		return "/"
	}

	return currentUser.HomeDir
}

func ABSPath(path string) string {
	absPath, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return ""
	}
	return absPath
}
