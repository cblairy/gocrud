package main

import (
	"golang.org/x/mod/modfile"
	"os"
)

func createFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func getModuleName() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return "", err
	}

	return modFile.Module.Mod.Path, nil
}
