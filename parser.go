package main

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Getting a Reader (html) and returning all the urls in the page in order as a slice. Urls are raw not parsed
func ParseUrlsFromHTML(r io.Reader) ([]string, error) {
	var urls []string
	doc, err := html.Parse(r)

	if err != nil {
		return nil, fmt.Errorf("parsing html %w", err)
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					urls = append(urls, a.Val)
					break
				}
			}
		}
	}

	return urls, nil
}

// <a href="<url>"> ....
func ParseUrlsFromHTMLMannual(html string) []string {
	var urls []string

	// outer loop
	for i := 0; i < len(html)-2; i++ {
		if html[i] == '<' && html[i+1] == 'a' && isTagEnd(html[i+2]) {
			// inner loop after finding a tag
			for j := i + 3; j < len(html); j++ {
				if html[j] == '>' {
					url := parseUrlFromAnchorTag(html[i : j+1]) // it will give <a>...href="">
					if url != "" {
						urls = append(urls, url)
					}
					break
				}
			}
		}
	}

	return urls
}

func isTagEnd(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '>' || b == '/'
}

func parseUrlFromAnchorTag(tag string) string {

	for i := 0; i < len(tag)-5; i++ {
		if tag[i] == 'h' && tag[i+1] == 'r' && tag[i+2] == 'e' && tag[i+3] == 'f' && tag[i+4] == '=' && tag[i+5] == '"' {
			// url := strings.Builder{}
			start := i + 6
			for j := start; j < len(tag); j++ {
				if tag[j] == '"' {
					return tag[start:j]

				}
			}
		}
	}

	return ""
}
