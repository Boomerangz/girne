package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeData(url string) (map[string]interface{}, error) {
	doc, err := GetParsedHtml(url)
	if doc == nil || err != nil {
		return nil, err
	}

	data := map[string]interface{}{"url": url, "price": []string{}}

	for _, field := range parseConfig.Fields {
		for _, selector := range field.Selector {
			if s := doc.Find(selector); s.Length() > 0 {
				s.Each(func(i int, s *goquery.Selection) {
					if field.Mode == MODE_ATTRIBUTE {
						for _, attribute := range field.Attribute {
							value, exists := s.Attr(attribute)
							if exists && len(value) > 0 {
								appendValue(data, field.Name, value)
							}
						}
					} else if field.Mode == MODE_TEXT {
						value := strings.TrimSpace(s.Text())
						if len(value) > 0 {
							appendValue(data, field.Name, value)
						}
					}
				})
			}
		}
	}
	return data, nil
}

func GetParsedHtml(url string) (*goquery.Document, error) {
	var resp *http.Response
	for triesCount := 0; triesCount < 10; triesCount++ {
		sem.Acquire()
		var err error
		resp, err = client.Get(url)
		if resp != nil && resp.StatusCode == 503 {
			r := rand.Intn(10)
			fmt.Println("503 sleep for ", r, "seconds")
			time.Sleep(time.Duration(r) * time.Second)
			continue
		}
		if err != nil {
			sem.Release()
			r := rand.Intn(10)
			fmt.Println(err, "sleep for ", r, "seconds")
			time.Sleep(time.Duration(r) * time.Second)
			continue
		} else if resp != nil && resp.StatusCode != 200 {
			sem.Release()
			return nil, errors.New(fmt.Sprintf("status code is %d", resp.StatusCode))
		}
		break
	}
	sem.Release()
	fmt.Println(url)
	doc, err := goquery.NewDocumentFromResponse(resp)
	resp.Body.Close()
	return doc, err
}

func appendValue(dictionary map[string]interface{}, key string, value string) {
	if _, ok := dictionary[key]; !ok {
		dictionary[key] = []string{value}
	} else {
		dictionary[key] = append(dictionary[key].([]string), value)
	}
}
