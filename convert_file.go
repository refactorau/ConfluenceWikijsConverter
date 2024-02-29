package main

import (
	"fmt"
	"os"
	"strings"
)

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
	fmt.Printf("Copy: %s/%s => %s/%s\n", source, fileName, destination, fileName)
	err = copyFile(source+"/"+fileName, destination+"/"+fileName)
	if err != nil {
		fmt.Printf("Copy Failed: %s/%s => %s/%s (%s)\n", source, fileName, destination, fileName, err)
	}

}
