// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// var urlTable = make(map[string][]string)
// var httpClient = &http.Client{Timeout: 10 * time.Second}
// var visitedUrls = VisitedURL{
// 	m: make(map[string]bool),
// }

// type SeedUrl struct {
// 	url   string
// 	depth int
// }

// type FetchedUrl struct {
// 	url     string
// 	results []string
// 	depth   int
// 	err     error
// }

// type VisitedURL struct {
// 	mu sync.Mutex
// 	m  map[string]bool
// }

// func (v *VisitedURL) CheckAndMark(url string) bool {
// 	v.mu.Lock()
// 	defer v.mu.Unlock()

// 	if v.m[url] {
// 		return true
// 	}
// 	v.m[url] = true
// 	return false
// }

// func worker(jobs <-chan SeedUrl, results chan<- FetchedUrl, workerWg *sync.WaitGroup) {
// 	defer workerWg.Done()

// 	for j := range jobs {
// 		foundURLs, err := crawl(j.url)
// 		results <- FetchedUrl{
// 			url:     j.url,
// 			results: foundURLs,
// 			depth:   j.depth,
// 			err:     err,
// 		}
// 	}
// }

// func crawl(url string) ([]string, error) {
// 	resp, err := httpClient.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("fetching url %s: %w", url, err)
// 	}
// 	defer resp.Body.Close()

// 	urls, err := ParseUrlsFromHTML(resp.Body, url)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return urls, nil
// }

// func main() {
// 	seedUrl := "https://info.cern.ch/" // first given url
// 	maxWorkers := 50
// 	jobs := make(chan SeedUrl)
// 	results := make(chan FetchedUrl, maxWorkers)
// 	workerWg := sync.WaitGroup{}
// 	depth := 5

// 	// workers pool
// 	for w := 1; w <= maxWorkers; w++ {
// 		workerWg.Add(1)
// 		go worker(jobs, results, &workerWg)
// 	}

// 	visitedUrls.CheckAndMark(seedUrl)
// 	pendingJobs := []SeedUrl{{url: seedUrl, depth: depth}}
// 	outstandingJobs := 1

// 	for outstandingJobs > 0 {
// 		var (
// 			nextJob SeedUrl
// 			jobSink chan SeedUrl
// 		)

// 		if len(pendingJobs) > 0 {
// 			nextJob = pendingJobs[0]
// 			jobSink = jobs
// 		}

// 		select {
// 		case jobSink <- nextJob:
// 			pendingJobs = pendingJobs[1:]
// 		case r := <-results:
// 			urlTable[r.url] = r.results

// 			if r.depth > 0 {
// 				for _, u := range r.results {
// 					if visitedUrls.CheckAndMark(u) {
// 						continue
// 					}
// 					pendingJobs = append(pendingJobs, SeedUrl{url: u, depth: r.depth - 1})
// 					outstandingJobs++
// 				}
// 			}

// 			if r.err != nil {
// 				fmt.Println(r.err)
// 			}

// 			outstandingJobs--
// 		}
// 	}

// 	close(jobs)
// 	workerWg.Wait()
// 	printUrls()
// }

// func printUrls() {
// 	for key, v := range urlTable {
// 		fmt.Printf("KEY %s ----\n ", key)
// 		for i := range v {
// 			fmt.Println(v[i])
// 		}
// 		fmt.Printf("\n\n")
// 	}
// 	fmt.Println("visited urls: ", visitedUrls.m, len(visitedUrls.m))
// }
