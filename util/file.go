package util

import (
	"bufio"
	"fmt"
	"lda/logging"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// FileExists checks if a file exists or not
func FileExists(filePath string) bool {
	if _, err := Fs.Stat(filePath); err == nil {
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
	file, err := Fs.Open(filePath)
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
	f, err := Fs.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}

// GetRepoNameFromConfig reads the .git/config file and extracts the repository name
func GetRepoNameFromConfig(path string) (string, error) {

	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	if err != nil && info != nil && !info.IsDir() {
		return "", fmt.Errorf("could not find .git directory: %w", err)
	}

	configPath := filepath.Join(gitPath, "config")
	file, err := Fs.Open(configPath)
	if err != nil {
		return "", fmt.Errorf("could not open .git/config: %w", err)
	}
	defer file.Close()

	var url string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Look for the URL in the origin section
		if strings.HasPrefix(line, "url =") {
			url = strings.TrimSpace(strings.TrimPrefix(line, "url ="))
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading .git/config: %w", err)
	}

	if url == "" {
		return "", fmt.Errorf("no origin URL found in .git/config")
	}

	url = strings.TrimSuffix(url, ".git") // Remove .git suffix if present
	parts := strings.Split(url, "/")

	// Extract repo name from URL
	return parts[len(parts)-1], nil
}
