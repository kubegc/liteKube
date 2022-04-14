package certificate

import "os"

// only all files exist return true, other return false
func Exists(files ...string) bool {
	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			return false
		}
	}
	return true
}
