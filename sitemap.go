package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetAllLinksFromRobots(url string) (chan string, error) {
	path := JoinUrl(url, "robots.txt")
	fmt.Println(path)
	data, err := HttpRequest(path)
	if err != nil {
		return nil, err
	}
	for _, st := range strings.Split(data, "\n") {
		if strings.Contains(st, "Sitemap") {
			sitemapAddres := strings.TrimSpace(strings.Replace(st, "Sitemap:", "", 1))
			if len(sitemapAddres) > 0 {
				ch := make(chan string, 50)
				go ParseSiteMap(sitemapAddres, ch, true)
				return ch, nil
			}
			break
		}
	}
	return nil, nil
}

func ParseSiteMap(sitemapUrl string, ch chan string, closeChannel bool) error {
	doc, err := goquery.NewDocument(sitemapUrl)
	if err != nil {
		return err
	}
	// Find the review items
	doc.Find("sitemap loc").Each(func(i int, s *goquery.Selection) {
		subMap := s.Text()
		ParseSiteMap(subMap, ch, false)
	})

	// Find the review items
	doc.Find("url loc").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		ch <- text
	})
	if closeChannel {
		close(ch)
	}
	return nil
}
