package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var urlTable = make(map[string][]string)
var visitedUrls = make(map[string]bool)

func crawl(url string, depth int, wg *sync.WaitGroup, mu *sync.Mutex) error {
	mu.Lock()
	if visitedUrls[url] {
		mu.Unlock()
		return nil
	}
	visitedUrls[url] = true // updating first as if err occured i want to keep it as visited marked
	mu.Unlock()

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
		for _, u := range urls {
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				crawl(u, depth-1, wg, mu)
			}(u)
		}
	}

	mu.Lock()
	urlTable[url] = urls
	mu.Unlock()

	return nil
}

func main() {
	url := "https://info.cern.ch/"
	// url := "https://news.ycombinator.com/"
	depth := 10
	var wg sync.WaitGroup
	var mu sync.Mutex

	err := crawl(url, depth, &wg, &mu)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
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
