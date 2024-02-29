package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func normaliseFilename(fileName string) string {
	re := regexp.MustCompile(`_\d+\.html$`)
	fileName = re.ReplaceAllString(fileName, ".html")
	return fileName
}

func normaliseDirectory(dir string) string {
	dir = strings.ReplaceAll(dir, " ", "-")
	dir = strings.ReplaceAll(dir, ":", "-")
	return dir
}

func fixFilesToMove(filesToMove *[]FileMove, root string) error {
	//fmt.Printf("%+v\n", filesToMove)

	for i, filemove := range *filesToMove {
		sourceFile := root + "/" + filemove.From
		destinationDir := filepath.Dir(filemove.To)
		destinationFilename := filepath.Base(filemove.To)

		data, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			fmt.Errorf("error reading file %s: %w", sourceFile, err)
			return err
		}

		// Parse the HTML content
		doc, err := html.Parse(bytes.NewReader(data))
		if err != nil {
			fmt.Errorf("error parsing HTML from file %s: %w", sourceFile, err)
			return err
		}

		// Attempt to find "breadcrumb-section" node
		breadcrumbSection, _ := findAndDetachNode(doc, "breadcrumb-section", false)
		//numDeep := 0
		if breadcrumbSection != nil {
			breadcrumbsNode, _ := findAndDetachNode(breadcrumbSection, "breadcrumbs", false)
			if breadcrumbsNode != nil {
				breadcrumbs := extractTextFromLI(breadcrumbsNode)
				for _, breadcrumb := range breadcrumbs {
					destinationDir = destinationDir + "/" + normaliseDirectory(breadcrumb)
				}
			}
		}

		if strings.HasPrefix(destinationDir, "./") {
			destinationDir = destinationDir[2:]
		}
		(*filesToMove)[i].To = destinationDir + "/" + normaliseFilename(destinationFilename)
	}

	return nil
}
