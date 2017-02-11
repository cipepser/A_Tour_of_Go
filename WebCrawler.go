package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

var fetched = struct {
	m map[string]error
	sync.Mutex
} {m: make(map[string]error)}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}
	
	// fetch済みか確認するためにlock
	fetched.Lock()
	
	// すでにurlがfetch済み
	if _, ok := fetched.m[url]; ok {
		fetched.Unlock()
		return
	}
	
	fetched.m[url] = fmt.Errorf("now fetching...")
	fetched.Unlock()

	body, urls, err := fetcher.Fetch(url)

	// fetchの結果を保存
	fetched.Lock()
	fetched.m[url] = err
	fetched.Unlock()

	// not foundならreturn
	if err != nil {
		fmt.Println(err)
		return
	}
	// foundなら結果を表示し、さらに下に潜る	
	fmt.Printf("found: %s %q\n", url, body)

	wait := new(sync.WaitGroup)
	for _, u := range urls {
		wait.Add(1)
		
		// Crawl url concurrently
		// uはfor文全体で使われる変数なので
		// 変数uをそのまま入れるとNG
		go func (url string) {
			Crawl(url, depth-1, fetcher)
			wait.Done()
		} (u)
	}
	wait.Wait()
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
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
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
