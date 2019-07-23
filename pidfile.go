// Package pidfile provides structure and helper functions to create and remove
// PID file. A PID file is usually a file used to store the process ID of a
// running process.
package pidfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PIDFile is a file used to store the process ID of a running process.
type PIDFile struct {
	path string
	pid  int
}

func checkPIDFileAlreadyExists(path string) error {
	if pidByte, err := ioutil.ReadFile(path); err == nil {
		pidString := strings.TrimSpace(string(pidByte))
		if pid, err := strconv.Atoi(pidString); err == nil {
			if processExists(pid) {
				return fmt.Errorf("pid file found, ensure docker is not running or delete %s", path)
			}
		}
	}
	return nil
}

// New creates a PIDfile using the specified path.
func New(path string) (*PIDFile, error) {
	if err := checkPIDFileAlreadyExists(path); err != nil {
		return nil, err
	}

	// Prepare dirs to ready file and ignore the error
	os.MkdirAll(filepath.Dir(path), os.FileMode(0755))

	pid := os.Getpid()

	if err := ioutil.WriteFile(path, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		return nil, err
	}

	return &PIDFile{path: path, pid: pid}, nil
}

// GetPID returns the pid of writen
func (file PIDFile) GetPID() int {
	return file.pid
}

// GetPath returns the path of pidfile
func (file PIDFile) GetPath() string {
	return file.path
}

// Remove removes the PIDFile.
func (file PIDFile) Remove() error {
	return os.Remove(file.path)
}
