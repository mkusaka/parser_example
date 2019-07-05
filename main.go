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

	statusResutls := []string{"url, statusCode"}

	for i, url := range targetURLs {
		if i >= 2 {
			break
		}
		targetURL := strings.Replace(url, replaceFromString, replaceToString, 1)
		resp, err := http.Get(targetURL)

		status := ""
		if err != nil {
			fmt.Println(err)
			status = err.Error()
		} else {
			status = strconv.Itoa(resp.StatusCode)
			resp.Body.Close()
		}

		statusResutls = append(statusResutls, targetURL+","+status)
	}

	file.Write(([]byte)(strings.Join(statusResutls, "\n")))
}
