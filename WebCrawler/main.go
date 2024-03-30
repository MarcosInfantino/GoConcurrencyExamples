package main

import (
	"fmt"
	"slices"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, parentWg *sync.WaitGroup) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if parentWg != nil {
		defer parentWg.Done()
	}

	actualWg := &sync.WaitGroup{}
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		if !urlCache.contains(u) {
			actualWg.Add(1)
			urlCache.append(u)
			go Crawl(u, depth-1, fetcher, actualWg)
		}
	}

	actualWg.Wait()

	return
}

func main() {
	Crawl("https://golang.org/", 4, fetcher, nil)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

type UrlCache struct {
	mutex sync.Mutex
	urls  []string
}

func (u *UrlCache) contains(url string) bool {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return slices.Contains(u.urls, url)
}

func (u *UrlCache) append(url string) {
	u.mutex.Lock()
	u.urls = append(u.urls, url)
	u.mutex.Unlock()
}

var urlCache = UrlCache{urls: []string{}}
