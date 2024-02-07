package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func listGoFiles(dirPath string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(d.Name()) == ".go" {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func saveFile(fileName, generatedCode string) {
	// create the file
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// write the generated code to the file
	if _, err = file.WriteString(generatedCode); err != nil {
		panic(err)
	}
}

func removeFile(fileName string) {
	// if exists, delete the file
	if _, err := os.Stat(fileName); err == nil {
		err = os.Remove(fileName)
		if err != nil {
			panic(err)
		}
	}
}

func isDirectory(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		log.Fatalf("failed to get file info for '%s'.\n%v", path, err)
	}
	return stat.IsDir()
}
