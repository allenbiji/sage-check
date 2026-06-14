package detect

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// looks for the filename to verify that it exists
func fileExists(filePath string) bool {
	file, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("There file is not accessible", err)
		return false
	}

	if err != nil {
		return false
	}

	return !file.IsDir()
}

// extract the env keys from the .env.example file
func extractEnvKeys(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("The file could not be opened", err)
		return nil
	}

	defer file.Close()

	var envKeys []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, _, found := strings.Cut(line, "=")
		if found {
			envKeys = append(envKeys, key)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file", err)
	}

	return envKeys
}
