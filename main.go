package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/mkusaka/sitemapparser"
	_ "github.com/motemen/go-loghttp/global"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal(`argument required
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
		return
	}
	defer file.Close()

	targetURLs, err := sitemapparser.Scheduler(indexSitemapURL)

	if err != nil {
		log.Fatal(err)
		return
	}

	type Resp struct {
		status string
		url    string
	}
	// statusResults := []string{"url, statusCode"}

	c := make(chan Resp, 10)
	for _, url := range targetURLs {
		go func() {
			targetURL := strings.Replace(url, replaceFromString, replaceToString, 1)
			resp, err := http.Get(targetURL)

			status := ""
			if err != nil {
				status = err.Error()
			} else {
				status = strconv.Itoa(resp.StatusCode)
				defer resp.Body.Close()
				c <- Resp{status, targetURL}
			}
		}()

		// statusResults = append(statusResults, statusResults+","+status)
	}

	for v := range c {
		fmt.Printf("status: %v, url: 5v", v.status, v.url)
	}
	// file.Write(([]byte)(strings.Join(statusResults, "\n")))
}
