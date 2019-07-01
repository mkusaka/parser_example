package main

import (
	"log"
	"os"
	"strings"

	"github.com/mkusaka/sitemapparser"
)

func main() {
	file, err := os.Create(`output.csv`)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	output, err := sitemapparser.Scheduler("some good gzipped sitemap url")

	if err != nil {
		log.Fatal(err)
	}

	file.Write(([]byte)(strings.Join(output, "\n")))
}
