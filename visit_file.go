package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func visitFile(root string, destination string, source string, target string, filesToMove []FileMove) error {

	sourceFile := root + "/" + source
	destFile := destination + "/" + target

	dir := filepath.Dir(destFile)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	if strings.HasSuffix(destFile, ".html") {
		err = convertHtml(root, destination, source, target, filesToMove)
		return err
	}

	// If we get here we dont need to convert anything, just copy the file
	fmt.Printf("Copy: %s => %s\n", sourceFile, destFile)
	err = copyFile(sourceFile, destFile)
	if err != nil {
		fmt.Printf("Copy Failed: %s => %s (%s)\n", sourceFile, destFile, err)
		return err
	}
	return nil
}
