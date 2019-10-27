package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"net/http"
	"strings"
)

const path = "http://anibin.blogspot.com/search/label/放送中"

type Resolution struct {
	Name string
	Res  string
}

var Resolutions []Resolution

func traverseHtml(n *html.Node, fn func(node *html.Node) (*html.Node, error)) (node *html.Node, err error) {
	node, err = fn(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if node != nil || err != nil {
			return
		}
		node, err = traverseHtml(c, fn)
	}
	return
}

func findTable(n *html.Node) (*html.Node, error) {
	if n.Type == html.ElementNode && n.Data == "table" {
		return n, nil
	}
	return nil, nil
}

var sb strings.Builder

func nodeToString(n *html.Node) (*html.Node, error) {
	if n.Type == html.TextNode {
		sb.WriteString(strings.TrimSpace(n.Data))
	}
	return nil, nil
}

var seenFirstRow bool

func addToResolutions(n *html.Node) (*html.Node, error) {
	if n.Type == html.ElementNode && n.Data == "tr" {
		if !seenFirstRow {
			seenFirstRow = true
		} else {
			sb.Reset()
			traverseHtml(n.FirstChild, nodeToString)
			name := strings.TrimSpace(sb.String())
			sb.Reset()
			traverseHtml(n.LastChild, nodeToString)
			res := strings.TrimSpace(sb.String())
			if len(name) > 0 && len(res) > 0 {
				resolution := Resolution{norm.NFKD.String(name), res}
				Resolutions = append(Resolutions, resolution)
			}
		}
	}
	return nil, nil
}

func GetResolutions() (err error) {
	fmt.Println("#### Beginning pull from Anibin for resolutions ####")
	resp, err := http.Get(path)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return
	}
	table, err := traverseHtml(doc, findTable)

	seenFirstRow = false
	traverseHtml(table, addToResolutions)

	fmt.Printf("INFO: Found %d entries on http://anibin.blogspot.com/search/label/%%E6%%94%%BE%%E9%%80%%81%%E4%%B8%%AD for current season\n",
		len(Resolutions))
	return
}
