package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func visit(n *html.Node) {
	if n.Type == html.ElementNode {
		fmt.Println(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		visit(c)
	}
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	visit(doc)
}
