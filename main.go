package main

import (
	"fmt"
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

	urls, err := ParseUrlsFromHTML(resp.Body)

	if err != nil {
		return err
	}

	fmt.Println(len(urls))
	fmt.Println(urls)

	return nil
}
func main() {
	err := crawl("https://news.ycombinator.com")
	if err != nil {
		log.Fatal(err)
	}
}
