package main

import (
	"fmt"
	"os"
	"path/filepath"
)

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
