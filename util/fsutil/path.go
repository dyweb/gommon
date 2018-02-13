package fsutil

import "os"

func Cwd() string {
	if cur, err := os.Getwd(); err != nil {
		log.Warnf("can't get current directory %v", err)
		return ""
	} else {
		return cur
	}
}

// Exists check if path exists and returns false regardless of detail error message
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// FileExists check if the path exists and it is a file, returns false regardless of detail error message
// i.e. it is possible the file exists but current reader doesn't have permission to read
func FileExists(path string) bool {
	if i, err := os.Stat(path); err != nil {
		return false
	} else {
		return !i.IsDir()
	}
}

// DirExists check if the path exists and it is a directory
func DirExists(path string) bool {
	if i, err := os.Stat(path); err != nil {
		return false
	} else {
		return i.IsDir()
	}
}
