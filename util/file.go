package util

import (
	"bufio"
	"lda/config"
	"lda/logging"
	"os"
	"os/user"
	"strconv"
	"strings"
)

// FileExists checks if a file exists or not
func FileExists(filePath string) bool {
	if _, err := config.Fs.Stat(filePath); err == nil {
		return true
	} else if !os.IsNotExist(err) {
		logging.Log.Err(err).Msg("Failed to check if file exists or not")
	}
	return false
}

// CreateDirAndChown creates a directory and changes its ownership
func CreateDirAndChown(dirPath string, perm os.FileMode, user *user.User) error {
	if err := os.MkdirAll(dirPath, perm); err != nil {
		return err
	}

	if user != nil {
		uid, err := strconv.Atoi(user.Uid)
		if err != nil {
			return err
		}

		gid, err := strconv.Atoi(user.Gid)
		if err != nil {
			return err
		}

		if err := os.Chown(dirPath, uid, gid); err != nil {
			return err
		}
	}

	return nil
}

// WriteFileAndChown writes content to a file and changes its ownership
func WriteFileAndChown(filePath string, content []byte, perm os.FileMode, user *user.User) error {
	if err := os.WriteFile(filePath, content, perm); err != nil {
		return err
	}

	if user != nil {

		uid, err := strconv.Atoi(user.Uid)
		if err != nil {
			return err
		}

		gid, err := strconv.Atoi(user.Gid)
		if err != nil {
			return err
		}

		if err := os.Chown(filePath, uid, gid); err != nil {
			return err
		}
	}

	return nil
}

// ChangeFileOwnership changes the ownership of a file
func ChangeFileOwnership(filePath string, user *user.User) error {
	if user == nil {
		return nil
	}

	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return err
	}

	if err := os.Chown(filePath, uid, gid); err != nil {
		return err
	}
	return nil
}

// IsScriptPresent checks if a script is already present in a file
func IsScriptPresent(filePath, script string) bool {
	file, err := config.Fs.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), script) {
			return true
		}
	}
	return false
}

// AppendToFile appends content to a file
func AppendToFile(filePath, content string) error {
	f, err := config.Fs.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}
