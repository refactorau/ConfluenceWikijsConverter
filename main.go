package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Get the file permissions of the source file
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	destinationFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func convertFile(source string, destination string, fileName string) {
	err := os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return
	}

	if strings.HasSuffix(fileName, ".html") {
		convertHtml(source, destination, fileName)
		return
	}

	// If we get here we dont need to convert anything, just copy the file
	err = copyFile(source+"/"+fileName, destination+"/"+fileName)
	if err != nil {
		fmt.Printf("Copy Failed: %s/%s => %s/%s (%s)\n", source, fileName, destination, fileName, err)
	}

}
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run program.go <directory_path> <destination_path>")
		return
	}
	root := os.Args[1] // get the directory path from the command line

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run program.go <directory_path> <destination_path>")
		return
	}
	destination := os.Args[2] // get the directory path from the command line

	err := filepath.Walk(root, func(fp string, fi os.FileInfo, err error) error {
		return visitFile(fp, fi, err, root, destination)
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}
}
