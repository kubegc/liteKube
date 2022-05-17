package utils

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"syscall"
)

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

func CreateDir(path string) {
	if !Exists(path) {
		os.MkdirAll(path, os.ModePerm)
	}
}

//func GetHomeDir() string {
//	if home, err := os.UserHomeDir(); err != nil {
//		return "/"
//	} else {
//		return home
//	}
//}

func GetHomeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		return ""
	}

	return currentUser.HomeDir
}

func CopyFile(src, dst string) error {

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if Exists(dst) {
		return fmt.Errorf("File %s already exists.", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return nil
	}
	defer destination.Close()

	buf := make([]byte, 1024)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

func Pwd() (string, error) {
	return os.Getwd()
}

func LockFile(path string) (*os.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB) // 加上排他锁，当遇到文件加锁的情况直接返回 Error
	if err != nil {
		return nil, fmt.Errorf("cannot flock file %s: %s", path, err)
	}
	return f, nil
}

func UnlockFile(file *os.File) error {
	defer file.Close()
	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}
