package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func convertHtml(root string, destination string, source string, target string, filesToMove []FileMove) error {
	sourceFile := root + "/" + source
	destFile := destination + "/" + target
	//fmt.Println("*******************************************", sourceFile, destFile)

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

	// Attempt to find and detach the "toc-macro" node
	findAndDetachNode(doc, "toc-macro", true)

	// Attempt to find and detach the "breadcrumb-section" node
	breadcrumbSection, detached := findAndDetachNode(doc, "breadcrumb-section", true)
	numDeep := -1
	if detached {
		// Now find the breadcrumbs node
		breadcrumbsNode, detached := findAndDetachNode(breadcrumbSection, "breadcrumbs", false)
		if detached {
			breadcrumbs := extractTextFromLI(breadcrumbsNode)
			for _, breadcrumb := range breadcrumbs {
				if numDeep == -1 {
					// This is the first link, so this is actually the home dir, so we dont need to add another directory
				} else {
					destination = destination + "/" + normaliseDirectory(breadcrumb)
				}
				numDeep = numDeep + 1
			}
		}
	}
	// Scan the entire HTML document, and every attribute. If the attribute starts with "attachments/" then update it to be "../attachments/"
	// ALSO normalise any "id" attributes while we are here.
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i, a := range n.Attr {
				// I cant believe I am doing this, but unless there are a LOT of files then its working fairly quickly
				for _, fileMove := range filesToMove {
					if strings.Contains(a.Val, fileMove.From) {
						newVal := strings.ReplaceAll(a.Val, fileMove.From, fileMove.To)
						if numDeep != -1 {
							newVal = strings.Repeat("../", numDeep) + newVal
						}
						//fmt.Println("=========", a.Val, "===========", newVal)
						n.Attr[i].Val = newVal
					}
				}

				if a.Key == "id" || strings.HasPrefix(a.Val, "#") {
					n.Attr[i].Val = normaliseIdentifier(a.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
	// After modifying the HTML document, write it back to a new file in the destination directory
	fmt.Printf("Convert: %s => %s\n", sourceFile, destFile)
	newFile, err := os.Create(destFile)
	if err != nil {
		fmt.Errorf("error creating file %s: %w", destFile, err)
		return err
	}
	defer newFile.Close()

	// Render the modified HTML document to the new file
	err = html.Render(newFile, doc)
	if err != nil {
		fmt.Errorf("error writing modified HTML to file %s: %w", destFile, err)
		return err
	}

	return nil
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
