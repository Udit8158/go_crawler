**This is a learning project to learn golang better**

## What we are building here?

- A crawler is like - we will give it a url. and it will return all the urls in that page. then it will go through those url again (here only same hostname urls) for few depth. So the final output will be something like a graph/tree of urls linked to the given page
- Example shape of output:
  `https://news.ycombinator.com (depth 0)
  ├── https://news.ycombinator.com/news (depth 1)
  │   ├── https://news.ycombinator.com/item?id=42 (depth 2)
  │   └── https://news.ycombinator.com/user?id=foo (depth 2)
  ├── https://news.ycombinator.com/newest (depth 1)
  │   └── ...
  └── https://news.ycombinator.com/ask (depth 1)
  └── ...`

## Goal

- learning about parsing html and url formats (which can be useful for building crawler which have actual utility)
- recursively doing this to build a graph/tree like data structure
- implement concurrency to boost speed (familiar with actual go overall)
- then will improve the overall implementaion (calude review) to make overall better product and learn better design

## Question came in my mind

- Why defer close request body?
  - TCP scoket, file descriptor, buffer in memory vs live network connection
- Why multiple go routines writting to an slice/map is bad but for channels it's good?

# Notes

- **Go routines with wg**
  - spawn go rotuines add wg.Add for each routine and call wg.Done for each. 
  - then before printing the final result add wg.Wait() so that it waits for each wg.Done() call 
  - now when the go routine function will write data in a data strcuture like array, maps we need to use wg.Mutex
  - why? because of **atomicity**. if multiple thing writes a ds in same time - they might read the value same and then update that value (so basically even after multiple writes the value only can change by one). and compiler doesn't allow that ofc.
  - also there are problems with maps for resizing and all (not clear to me yet)
  - that's the reason we have to use mutex.Lock before reading/writting a ds and then unlock that. that's how atomicity will remain. and we can test that from go run -race (command to test *race condition - where mutiple threads write to same location*)
