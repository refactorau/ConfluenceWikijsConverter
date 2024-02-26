package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func normaliseFilename(fileName string) string {
	re := regexp.MustCompile(`_\d+.html$`)
	fileName = re.ReplaceAllString(fileName, "") + ".html"
	return fileName
}

func normaliseDirectory(dir string) string {
	dir = strings.ReplaceAll(dir, " ", "-")
	dir = strings.ReplaceAll(dir, ":", "-")
	return dir
}

func convertHtml(source string, destination string, fileName string) {
	filePath := source + "/" + fileName

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Errorf("error reading file %s: %w", filePath, err)
		return
	}

	// Parse the HTML content
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		fmt.Errorf("error parsing HTML from file %s: %w", filePath, err)
		return
	}

	// Attempt to find and detach the "main-header" node
	mainHeader, detached := findAndDetachNode(doc, "main-header")
	if detached {
		// Now find the breadcrumbs node
		breadcrumbsNode, detached := findAndDetachNode(mainHeader, "breadcrumbs")
		if detached {
			breadcrumbs := extractTextFromLI(breadcrumbsNode)
			for _, breadcrumb := range breadcrumbs {
				destination = destination + "/" + normaliseDirectory(breadcrumb)
			}
		}
	}

	// After modifying the HTML document, write it back to a new file in the destination directory
	os.MkdirAll(destination, os.ModePerm)
	fileName = normaliseFilename(fileName)
	newFilePath := destination + "/" + fileName
	newFile, err := os.Create(newFilePath)
	if err != nil {
		fmt.Errorf("error creating file %s: %w", newFilePath, err)
		return
	}
	defer newFile.Close()

	// Render the modified HTML document to the new file
	err = html.Render(newFile, doc)
	if err != nil {
		fmt.Errorf("error writing modified HTML to file %s: %w", newFilePath, err)
		return
	}

	return
}

// extractTextFromLI extracts text from each <li> element within the provided <ol> node.
func extractTextFromLI(n *html.Node) []string {
	var items []string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "span" {
					for d := c.FirstChild; d != nil; d = d.NextSibling {
						if d.Type == html.ElementNode && d.Data == "a" {
							for e := d.FirstChild; e != nil; e = e.NextSibling {
								if e.Type == html.TextNode {
									items = append(items, e.Data)
								}
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(n)
	return items
}

// findAndDetachNode searches and detaches a node by id.
// It returns the detached node and a boolean indicating if the operation was successful.
func findAndDetachNode(n *html.Node, id string) (*html.Node, bool) {
	var parent *html.Node
	var nodeToDetach *html.Node

	var traverse func(*html.Node) bool
	traverse = func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					nodeToDetach = n
					return true
				}
			}
		}
		parent = n
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if traverse(c) {
				return true
			}
		}
		parent = n.Parent // Reset parent as we backtrack
		return false
	}

	found := traverse(n)
	if found && parent != nil && nodeToDetach != nil {
		// Detach nodeToDetach from its parent
		if nodeToDetach.PrevSibling != nil {
			nodeToDetach.PrevSibling.NextSibling = nodeToDetach.NextSibling
		} else {
			parent.FirstChild = nodeToDetach.NextSibling
		}
		if nodeToDetach.NextSibling != nil {
			nodeToDetach.NextSibling.PrevSibling = nodeToDetach.PrevSibling
		}
		nodeToDetach.Parent = nil
		return nodeToDetach, true
	}
	return nil, false
}
