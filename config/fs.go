package config

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/samber/lo"
)

func configFilePathFromEnv(env string) string {
	// Get the directory of this file.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Unable to get cwd")
	}

	// Open the directory with the env configs.
	dirname := filepath.Dir(filename)
	dir, err := os.Open(dirname + "/env")
	if err != nil {
		log.Fatalf("Error opening config directory: %+v\n", err)
	}

	// Get the list of env config files.
	files, err := dir.ReadDir(0)
	if err != nil {
		log.Fatalf("Error reading config directory: %+v\n", err)
	}

	// Search the config files for one matching the current app environment.
	dirEntry, found := lo.Find(files, func(f fs.DirEntry) bool {
		return fileNameWithoutExt(f) == env
	})

	if !found {
		log.Fatalf("Could not find config file for environment %s\n", env)
	}

	// Return the config file path for the current app environment.
	return dirname + "/env/" + dirEntry.Name()
}

func fileNameWithoutExt(file fs.DirEntry) string {
	fileName := file.Name()

	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
