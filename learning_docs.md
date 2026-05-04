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

  - Problem: right now if i found 1000 url in a page it will spawn 1000 go routines and then more for nested. in scale it will cause problem
  - "for limit thing right now we are spawning a go routine for each url and it goes like
    first seed url -> main thread then each children url -> go routines -> for each children
     url's url -> go routuine (if depth allows). and it's like first children url the
    directly his children urls and then after depth finishes, then 2nd actual url. that's
    the design right now right?" - i was thinking that but it's not the design. fuck.


- Channel concurrency
  - buffered channel means having a limited cap (not zero), channel like that will not block the routine until the buffered limit hit. unlike a non buffered channel will be always blocking (until someone one recieving value) as it has 0 cap
  - v,ok <- ch here ok signify is ch is closed ( close*=(ch) for closing). after closing sending to a ch will cause panic, so it's sender's job (the go routine who send to the channel) to close it. recieving from a ch (close) will give zero values and can be checked like that
  - also there can be another channel we can use to synchronise. let's say something like we want to block main routine till another child routine needs to finish. for this wg.wait is also a good idea.
  - deadlock will only occur when all the go routuines are blocked and there is possibility to unblocked in near future. also there can be a thing where main is blocked indefinitely and other go routines are running but not doing something related to unblock the main - there it is called live blocking or goroutine leak. *so healthy program is when we are blocking something definitely. every go routine should have clear escape path*

 - In this proj i think i am being too much dependent on ai for design stuff, which is bad, that's causing me problem. i am being confused more
 - i think each thread or go routine should have clear responsibilites/job if it try to do mix and match stuff, things will likely break
 - if you want to send data from one channel to another channel, you can, but think carefully in scale that might occur deadlock. so alterantive approach would be write the data first into the memory and read from that in seperate channel. ofc you can ensure the atomicity of the data by mutex
