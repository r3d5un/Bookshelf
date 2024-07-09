package testing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get current working directory: %w", err)
	}

	markerFile := ".git"

	for {
		if _, err := os.Stat(filepath.Join(cwd, markerFile)); err == nil {
			return cwd, nil
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			return "", fmt.Errorf("project root not found")
		}
		cwd = parent
	}
}

func ListUpMigrationScrips(dirPath string) (migrations []string, err error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			migrations = append(migrations, fmt.Sprintf("%s/%s", dirPath, file.Name()))
		}
	}

	return migrations, nil
}
