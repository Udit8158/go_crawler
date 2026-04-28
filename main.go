package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	parser "example.com/m/lib"
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

	urls := parser.ParseUrlsFromHTML(html)
	fmt.Println(len(urls))

	return nil
}
func main() {
	err := crawl("https://news.ycombinator.com")
	if err != nil {
		log.Fatal(err)
	}
}
