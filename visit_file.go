package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func visitFile(fp string, fi os.FileInfo, err error, source string, destination string) error {
	if err != nil {
		fmt.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}
	if !!fi.IsDir() {
		return nil // not a file. ignore.
	}

	dir := filepath.Dir(fp)       // get the directory of the file
	fileName := filepath.Base(fp) // get the filename

	if strings.HasPrefix(dir, source) {
		dest_dir := destination + dir[len(source):]
		convertFile(dir, dest_dir, fileName)
	}
	return nil
}
