package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func GetHost(urlstr string) string {
	parsed, err := url.Parse(urlstr)
	if err == nil {
		return parsed.Host
	}
	return ""
}

func JoinUrl(siteUrl, path string) string {
	u, err := url.Parse(path)
	if err != nil {
		log.Fatal(err)
	}
	base, err := url.Parse(siteUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u, base)
	return base.ResolveReference(u).String()
}

func HttpRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}
