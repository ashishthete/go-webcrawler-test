package main

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

type Parser interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Parse(r io.Reader) []string
}

type HtmlUrlParser struct {
	Filter func(url string) bool
	Urls   []string
}

func (p *HtmlUrlParser) extractLinks(doc *goquery.Document) []string {
	foundUrls := []string{}
	if doc != nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			res, _ := s.Attr("href")
			foundUrls = append(foundUrls, res)
		})
		return foundUrls
	}
	return foundUrls
}

func (p *HtmlUrlParser) resolveRelative(hrefs []string) []string {
	urls := []string{}
	for _, href := range hrefs {
		if href == "" || !p.Filter(href) {
			continue
		}
		urls = append(urls, href)
	}
	return urls
}

func (p *HtmlUrlParser) Parse(body io.Reader) []string {
	doc, _ := goquery.NewDocumentFromReader(body)
	links := p.extractLinks(doc)
	return p.resolveRelative(links)
}
