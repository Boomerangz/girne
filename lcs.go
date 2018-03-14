package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func LCS(s1 string, s2 string) string {
	var m = make([][]int, 1+len(s1))
	for i := 0; i < len(m); i++ {
		m[i] = make([]int, 1+len(s2))
	}
	longest := 0
	x_longest := 0
	for x := 1; x < 1+len(s1); x++ {
		for y := 1; y < 1+len(s2); y++ {
			if s1[x-1] == s2[y-1] {
				m[x][y] = m[x-1][y-1] + 1
				if m[x][y] > longest {
					longest = m[x][y]
					x_longest = x
				}
			}
		}
	}
	return s1[x_longest-longest : x_longest]
}

func findLCS(filename string, keys []string) map[string]string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	lcs := map[string]string{}
	for scanner.Scan() {
		text := scanner.Text()
		var data map[string][]string
		err := json.Unmarshal([]byte(text), &data)
		if err != nil {
			continue
		}

		for _, k := range keys {
			if dataValue, ok := data[k]; ok {
				for _, str := range dataValue {
					if _, ok := lcs[k]; ok {
						lcs[k] = LCS(lcs[k], str)
					} else {
						lcs[k] = str
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()
	return lcs
}

func processFileForLCS(filename string, ouputfile string, lcs map[string]string) {
	postFilename := filename + ".post"
	postFile, err := os.Create(postFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		var data map[string][]string
		err := json.Unmarshal([]byte(text), &data)
		if err != nil {
			continue
		}

		for k, lcString := range lcs {
			if dataValue, ok := data[k]; ok {
				for idx := range dataValue {
					data[k][idx] = strings.Replace(data[k][idx], lcString, "", 1)
				}
			}
		}

		jsonData, err := json.Marshal(data)
		if err == nil {
			postFile.Write(jsonData)
			postFile.WriteString("\n")
		} else {
			fmt.Println(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()
	postFile.Close()
}
