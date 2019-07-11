package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/mkusaka/sitemapparser"
	_ "github.com/motemen/go-loghttp/global"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal(`argument requierd
1st argument: index sitemap url
2nd argument: replace from collected site url string
3nd argument: replace to collected site url string`)
		return
	}
	indexSitemapURL := os.Args[1]

	replaceFromString := ""
	if len(os.Args) >= 3 {
		replaceFromString = os.Args[2]
	}

	replaceToString := ""
	if len(os.Args) >= 4 {
		replaceToString = os.Args[3]
	}

	file, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	targetURLs, err := sitemapparser.Scheduler(indexSitemapURL)

	if err != nil {
		log.Fatal(err)
	}

	statusResults := []string{"url, statusCode"}

	// parallel access count
	maxConnection := make(chan bool, 10)
	wg := &sync.WaitGroup{}

	for i, url := range targetURLs {
		wg.Add(1)
		maxConnection <- true

		go func() {
			defer wg.Done()

			targetURL := strings.Replace(url, replaceFromString, replaceToString, 1)
			resp, err := http.Get(targetURL)

			status := ""
			if err != nil {
				status = err.Error()
			} else {
				status = strconv.Itoa(resp.StatusCode)
				resp.Body.Close()
			}
			statusResults = append(statusResults, targetURL+","+status)
			fmt.Println("current: " + strconv.Itoa(i) + "/" + strconv.Itoa(len(targetURLs)))
		}()
	}

	file.Write(([]byte)(strings.Join(statusResults, "\n")))
}
