package main

import (
	"fmt"
	"io"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Getting a Reader (html) and returning all the urls in the page in order as a slice. Urls are raw not parsed
func ParseUrlsFromHTML(r io.Reader, u string) ([]string, error) {
	var urls []string
	baseUrl, err := url.Parse(u)

	if err != nil {
		return nil, fmt.Errorf("parsing url %s: %w", u, err)
	}

	doc, err := html.Parse(r)

	if err != nil {
		return nil, fmt.Errorf("parsing html %w", err)
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					parsedUrl, err := url.Parse(a.Val)
					if err != nil {
						continue
					}

					if !parsedUrl.IsAbs() {
						// realive url
						parsedUrl = baseUrl.ResolveReference(parsedUrl)
					}

					if parsedUrl.Host == baseUrl.Host && (parsedUrl.Scheme == "https" || parsedUrl.Scheme == "http") {
						parsedUrl.Fragment = ""
						urls = append(urls, parsedUrl.String())
					}
					break
				}
			}
		}
	}

	return urls, nil
}
