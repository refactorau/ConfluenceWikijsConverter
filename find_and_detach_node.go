package main

import "golang.org/x/net/html"

// findAndDetachNode searches and detaches a node by id.
// It returns the detached node and a boolean indicating if the operation was successful.
func findAndDetachNode(n *html.Node, id string, doDetach bool) (*html.Node, bool) {
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
	}

	// Start the search and detach process
	if nodeToDetach != nil {
		return nodeToDetach, true
	}
	return nil, false
}
