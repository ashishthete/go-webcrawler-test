package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

type UrlData struct {
	childs []string
	depth  int
}
type Crawler struct {
	crawled map[string]UrlData
	mux     sync.Mutex
}

func New() *Crawler {
	return &Crawler{
		crawled: make(map[string]UrlData),
	}
}

func (c *Crawler) Print() {
	fmt.Println(c.crawled)
}

func (c *Crawler) visited(url string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	_, ok := c.crawled[url]
	return ok
}

func (c *Crawler) AddRoute(url string, urls []string, depth int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.crawled[url] = UrlData{childs: urls, depth: depth}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func (c *Crawler) Crawl(url string, fetcher Fetcher, parser Parser, depth, expectedDepth int) {
	var wg sync.WaitGroup
	if depth > expectedDepth {
		return
	}

	if c.visited(url) {
		return
	}
	body, err := fetcher.Fetch(url)
	if err != nil {
		// fmt.Println(err)
		return
	}

	urls := parser.Parse(strings.NewReader(string(body)))

	c.AddRoute(url, urls, depth)

	for _, u := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			c.Crawl(u, fetcher, parser, depth+1, expectedDepth)
		}(u)
	}
	wg.Wait()
}

func getHost(u string) (string, error) {
	parsed, err := url.Parse(u)
	return parsed.Host, err
}

func main() {
	crawler := New()
	fetcher := HttpFetcher{}
	// startUrl := "https://www.barclays.in/"
	startUrl := "https://www.google.com/"
	expectedDepth := 2 // 0 means unlimited

	host, err := getHost(startUrl)
	if err != nil {
		panic("invalid start url")
	}

	hostParts := len(strings.Split(host, "."))

	filter := func(url string) bool {
		urlHost, err := getHost(url)
		if err != nil {
			return false
		}
		urlTokens := strings.Split(urlHost, ".")
		index := len(urlTokens) - hostParts
		if index < 0 {
			return false
		}

		return strings.Join(urlTokens[index:], ".") == host
	}
	parser := &HtmlUrlParser{Filter: filter}
	crawler.Crawl(startUrl, fetcher, parser, 1, expectedDepth)

	crawler.Print()
}
