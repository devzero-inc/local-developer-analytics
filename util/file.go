package util

import (
	"bufio"
	"lda/logging"
	"os"
	"strings"
)

// FileExists checks if a file exists or not
func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if !os.IsNotExist(err) {
		logging.Log.Err(err).Msg("Failed to check if file exists or not")
	}
	return false
}

// IsScriptPresent checks if a script is already present in a file
func IsScriptPresent(filePath, script string) bool {
	file, err := os.Open(filePath)
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
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}
