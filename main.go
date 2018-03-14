package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/abiosoft/semaphore"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	link           = kingpin.Flag("url", "Website url.").Required().String()
	threadsCount   = kingpin.Flag("threads_count", "Threads count").Default("10").Int()
	downloadsCount = kingpin.Flag("downloads_count", "Downloads count").Default("10").Int()
	timeoutSeconds = kingpin.Flag("timeout", "Maximum timeout").Default("15").Int()
	outputFile     = kingpin.Flag("output", "Output file.").String()

	sem     *semaphore.Semaphore
	timeout time.Duration
	client  http.Client
)

func main() {
	kingpin.Parse()

	sem = semaphore.New(*downloadsCount)
	timeout = time.Duration(*timeoutSeconds) * time.Second
	client = http.Client{
		Timeout: timeout,
	}

	start := time.Now()

	filename := GetHost(*link) + ".json"
	if outputFile != nil && len(*outputFile) > 0 {
		filename = *outputFile
	}
	ch, err := GetAllLinksFromRobots(*link)
	if err != nil {
		return
	}
	result := make(chan map[string]interface{}, 50)
	go ParseAllPrices(ch, result)
	WriteToFileFromChannel(result, filename)
	PostProcess(filename)

	fmt.Println(time.Since(start))
}

func PostProcess(filename string) {
	filenamePost := filename + ".post"
	keys := []string{"title", "name", "og:title", "description"}
	lcs := findLCS(filename, keys)
	processFileForLCS(filename, filenamePost, lcs)
	os.Remove(filename)
	os.Rename(filenamePost, filename)
}

func WriteToFileFromChannel(outputChannel chan map[string]interface{}, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		data, more := <-outputChannel
		if more {
			jsonData, err := json.Marshal(data)
			if err == nil {
				f.Write(jsonData)
				f.WriteString("\n")
			} else {
				fmt.Println(err)
			}
		} else {
			break
		}
	}
	f.Close()
}

func ParseAllPricesInner(urls chan string, results chan map[string]interface{}) {
	r := rand.Intn(*threadsCount * 2)
	fmt.Println("sleep for ", r, "seconds")
	time.Sleep(time.Duration(r) * time.Second)
	for {
		link, more := <-urls
		if more {
			data, err := ScrapeData(link)
			if err != nil {
				fmt.Printf("link %s error %s\n", link, err)
			} else if data == nil {
				fmt.Printf("link %s result nil\n", link)
			} else {
				if _, ok := data["price"].([]string); ok && len(data["price"].([]string)) > 0 {
					results <- data
				} else if _, ok := data["alt_price"]; ok {
					results <- data
				}
			}
		} else {
			fmt.Println("finished ParseAllPricesInner")
			break
		}
	}
	close(results)
}

func ParseAllPrices(urls chan string, results chan map[string]interface{}) {
	channels := []chan map[string]interface{}{}
	for i := 0; i < *threadsCount; i++ {
		ch := make(chan map[string]interface{}, 50)
		channels = append(channels, ch)
		go ParseAllPricesInner(urls, ch)
	}

	for {
		dataCollected := false
		for _, ch := range channels {
			data, more := <-ch
			if more {
				dataCollected = true
				results <- data
			}
		}
		if !dataCollected {
			break
		}
	}
	close(results)
}
