package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	url   string
	depth int
}
type Result struct {
	parent string
	urls   []string
	depth  int
}
type Visited struct {
	mu     sync.Mutex
	urlMap map[string]bool
}

func (v *Visited) checkAndMark(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.urlMap[url] {
		return true
	}
	v.urlMap[url] = true
	return false
}

var httpClient = &http.Client{Timeout: 10 * time.Second}
var urlsTable = make(map[string][]string)
var visitedUrls = Visited{urlMap: make(map[string]bool)}

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		foundUrlsInPage, _ := crawl(j.url)
		results <- Result{parent: j.url, urls: foundUrlsInPage, depth: j.depth - 1}
	}
}

func main() {
	startTime := time.Now()

	maxWorkers := 100000
	jobs := make(chan Job)
	results := make(chan Result, maxWorkers)
	depth := 30
	workerWg := sync.WaitGroup{}
	// seedUrl := "https://info.cern.ch/"
	seedUrl := "https://news.ycombinator.com/"

	jobQueue := []Job{}
	jobQueue = append(jobQueue, Job{url: seedUrl, depth: depth})
	visitedUrls.checkAndMark(seedUrl)
	outstandignJobs := 1

	// spawn workder N
	for w := 1; w <= maxWorkers; w++ {
		workerWg.Add(1)
		go worker(jobs, results, &workerWg)
	}

	// processing that results -> to create more jobs (for depth) and printing data or collecting data in map
	go func() {
		for {
			var (
				nextJob Job
				jobSink chan Job
			)

			if len(jobQueue) == 0 && outstandignJobs == 0 {
				close(jobs)
				return // break the loop
			}

			if len(jobQueue) > 0 {
				nextJob = jobQueue[0]
				jobSink = jobs
			}

			select {
			case r := <-results:
				urlsTable[r.parent] = r.urls
				outstandignJobs-- // one result coming means on job done
				if r.depth > 0 {
					// outstandignJobs++
					// outstandignJobs += len(r.urls)
					for _, u := range r.urls {
						if !visitedUrls.checkAndMark(u) {
							outstandignJobs++
							jobQueue = append(jobQueue, Job{url: u, depth: r.depth})
						}
					}
				}
			case jobSink <- nextJob:
				jobQueue = jobQueue[1:]
			}
		}
	}()

	workerWg.Wait()
	printUrls()
	close(results)

	endTime := time.Since(startTime)
	fmt.Println("total time taken: ", endTime)
}

func crawl(url string) ([]string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching url %s: %w", url, err)
	}
	defer resp.Body.Close()

	urls, err := ParseUrlsFromHTML(resp.Body, url)

	if err != nil {
		return nil, err
	}
	return urls, nil
}

func printUrls() {
	for key, v := range urlsTable {
		fmt.Printf("KEY %s ----\n ", key)
		for i := range v {
			fmt.Println(v[i])
		}
		fmt.Printf("\n\n")
	}
	fmt.Println("visited urls: ", visitedUrls.urlMap, len(visitedUrls.urlMap))
}
