package global

import (
	"fmt"
	"io/ioutil"
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

func CopyFile(source string, destination string) error {
	fmt.Print(source, "\n", destination, "\n")
	sourceBytes, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(destination, sourceBytes, os.FileMode(0644)); err != nil {
		return err
	}
	return nil
}

// only all files exist return true, other return false
func Exists(files ...string) bool {
	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			return false
		}
	}
	return true
}

func NotExists(files ...string) bool {
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			return false
		}
	}
	return true
}
