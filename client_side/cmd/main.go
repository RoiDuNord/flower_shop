package main

import (
	"fmt"
	"scraper"
)

func main() {
	flowers := scraper.Scrape()
	fmt.Println(flowers)
}
