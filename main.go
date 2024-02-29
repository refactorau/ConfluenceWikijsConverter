package main

import (
	"fmt"
	"os"
	"strings"
)

type FileMove struct {
	From string
	To   string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run program.go <directory_path> <destination_path>")
		return
	}
	root := os.Args[1] // get the directory path from the command line
	if strings.HasSuffix(root, "/") {
		root = root[:len(root)-1]
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run program.go <directory_path> <destination_path>")
		return
	}
	destination := os.Args[2] // get the directory path from the command line
	if strings.HasSuffix(destination, "/") {
		destination = destination[:len(destination)-1]
	}

	// First we just walk the directory tree and list all of the files we have.
	filesToMove, err := listAllFiles(root)
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}

	// Now the filesToMove are a list of files. Next we might want to update the destination of the files...
	fixFilesToMove(&filesToMove, root)

	// Now we can either copy the files to the destination, or update them and write them out
	for _, f := range filesToMove {
		visitFile(root, destination, f.From, f.To, filesToMove)
	}
}
