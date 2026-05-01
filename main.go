package main

import (
	"fmt"
	"log"
	"net/http"
)

var urlTable = make(map[string][]string)
var visitedUrls = make(map[string]bool)

func crawl(url string, depth int) error {
	if visitedUrls[url] {
		return nil
	}
	visitedUrls[url] = true // updating first as if err occured i want to keep it as visited marked

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("fetching url %s: %w", url, err)
	}

	// close the reader to prevent from resource leaks
	defer resp.Body.Close()

	urls, err := ParseUrlsFromHTML(resp.Body, url)
	// we are not handling the base url when our get request is gonna redirect
	// there we have to use resp.Request.URL (that's the redirected url)
	// we have to use that as a base url or else same seed hostname condition will break

	if err != nil {
		return err
	}

	if depth > 0 {
		for i := range urls {
			crawl(urls[i], depth-1)
		}
	}

	urlTable[url] = urls
	return nil
}

func main() {
	url := "https://info.cern.ch/"
	// url := "https://news.ycombinator.com/"
	depth := 1
	err := crawl(url, depth)
	if err != nil {
		log.Fatal(err)
	}
	printUrls()
}

func printUrls() {
	for key, v := range urlTable {
		fmt.Printf("KEY %s ----\n ", key)
		for i := range v {
			fmt.Println(v[i])
		}
		fmt.Printf("\n\n")
	}
	fmt.Println(visitedUrls, len(visitedUrls))
}
