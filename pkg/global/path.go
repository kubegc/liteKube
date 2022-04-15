package global

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

var HomePath string = GetHomeDir()
var LocalPath string = GetCurrentDir()

func GetCurrentDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return pwd
}

func GetHomeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		return ""
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

// replace ~, $HOME in path to ABS Path
func ReplaceHome(raw []byte) []byte {
	re := regexp.MustCompile(`(,|"|\s)(~|\$HOME)/`)
	homePath := HomePath
	return re.ReplaceAll(raw, []byte(fmt.Sprintf("${1}%s/", strings.TrimRight(homePath, "/"))))
}

// replace . in path to ABS Path
func ReplaceCurrent(raw []byte) []byte {
	re := regexp.MustCompile(`(,|"|\s)(\.)/`)
	return re.ReplaceAll(raw, []byte(fmt.Sprintf("${1}%s/", strings.TrimRight(LocalPath, "/"))))
}
