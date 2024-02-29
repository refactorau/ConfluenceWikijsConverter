package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func listAllFiles(root string) ([]FileMove, error) {
	filesToMove := []FileMove{}

	err := filepath.Walk(root, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if !!fi.IsDir() {
			return nil // not a file. ignore.
		}
		dir := filepath.Dir(fp) // get the directory of the file
		if strings.HasPrefix(dir, root) {
			filesToMove = append(filesToMove, FileMove{fp[len(root)+1:], fp[len(root)+1:]})
		}
		return nil
	})
	if err != nil {
		return filesToMove, err
	}

	return filesToMove, nil
}
