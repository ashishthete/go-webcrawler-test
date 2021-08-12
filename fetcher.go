package main

import (
	"io/ioutil"
	"net/http"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) ([]byte, error)
}

// HttpFetcher is Fetcher that returns canned results.
type HttpFetcher struct{}

func (f HttpFetcher) Fetch(url string) ([]byte, error) {
	resp, err := getRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func getRequest(url string) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
