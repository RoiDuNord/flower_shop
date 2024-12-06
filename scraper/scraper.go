package scraper

import (
	"log"

	"github.com/gocolly/colly"
)

func Scrape() []string {
	var flowerArray = make([]string, 0)

	c := colly.NewCollector()

	c.OnHTML("strong", func(e *colly.HTMLElement) {
		text := e.Text

		var name string
		for _, letter := range text {
			if isCyrillic(letter) {
				name += string(letter)
			}
		}
		flowerArray = append(flowerArray, name)
	})

	err := c.Visit("https://flowerbiz.ru/bloger/top-20-zhivykh-tsvetov-dlya-sozdaniya-krasivykh-buketov/")
	if err != nil {
		log.Fatal(err)
	}

	return flowerArray
}

func isCyrillic(letter rune) (b bool) {
	if (letter >= 'А' && letter <= 'я') || letter == 'ё' || letter == 'Ё' {
		b = true
	}
	return b
}
