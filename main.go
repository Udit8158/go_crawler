package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func crawl(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("fetching url %s: %w", url, err)
	}

	// close the reader to prevent from resource leaks
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("reading responese body: %w", err)
	}

	html := string(data)

	tags := parseAnchorTagsFromHtml(html)
	fmt.Println(len(parseUrlsFromAnchorTags(tags)))

	return nil
}

func isTagEnd(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '>' || b == '/'
}

// <a href="<url>"> ....
func parseAnchorTagsFromHtml(html string) []string {
	var tags []string

	// outer loop
	for i := 0; i < len(html)-2; i++ {
		if html[i] == '<' && html[i+1] == 'a' && isTagEnd(html[i+2]) {
			// inner loop after finding a tag
			for j := i + 3; j < len(html); j++ {
				if html[j] == '>' {
					tags = append(tags, html[i:j])
					break
				}
			}
		}
	}

	return tags
}

func parseUrlsFromAnchorTags(tags []string) []string {
	var urls []string

	for _, tag := range tags {
		// println(tag)

		for i := 0; i < len(tag)-5; i++ {
			if tag[i] == 'h' && tag[i+1] == 'r' && tag[i+2] == 'e' && tag[i+3] == 'f' && tag[i+4] == '=' && tag[i+5] == '"' {
				// url := strings.Builder{}
				start := i + 6
				for j := start; j < len(tag); j++ {
					if tag[j] == '"' {
						urls = append(urls, tag[start:j])
						break
					}
				}

				// breaking from outer loop after finding the href
				break
			}
		}
	}

	return urls
}

func main() {
	err := crawl("https://news.ycombinator.com")
	if err != nil {
		log.Fatal(err)
	}
}
