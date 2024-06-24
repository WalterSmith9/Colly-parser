package main

import (
	"fmt" //formatted I/O
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly" //scraping framework
)

func main() {

	c := colly.NewCollector(colly.AllowedDomains("www.ozon.ru"))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Link of the page:", r.URL)
	})

	//getting description
	c.OnHTML("div.mm4_27", func(h *colly.HTMLElement) {
		name := h.ChildText("h1.mm3_27.tsHeadline550Medium")
		fmt.Println("Name: ", name)
	})

	//getting bread crumbs (path)
	c.OnHTML("ol.fe4_10", func(h *colly.HTMLElement) {
		var breadCrumbs string
		h.ForEach(".ah6.fe5_10", func(i int, h *colly.HTMLElement) {
			breadCrumbs += h.ChildText("span")
			breadCrumbs += "--"
		})
		breadCrumbs = strings.Trim(breadCrumbs, "-")
		fmt.Println(breadCrumbs)
	})

	//getting characteristics
	c.OnHTML("div[class=mm4_27]", func(h *colly.HTMLElement) {
		e, err := h.DOM.Html()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(e)

		h.ForEach("dl.p0k_27", func(i int, h *colly.HTMLElement) {
			key := h.ChildText("span.k0p_27")
			val := h.ChildText("dd.kp0_27")
			fmt.Printf("%s - %s/n", key, val)
		})

	})

	//testing: writing a DOM into .html file
	file, err := os.Create("pages/phone.html")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()

	c.OnHTML("html", func(h *colly.HTMLElement) {
		h.DOM.Each(func(i int, s *goquery.Selection) {
			e, err := s.Html()
			if err != nil {
				log.Println(err)
			}
			file.WriteString(e)
			//fmt.Println(e)
		})

	})

	//setting url for parsing
	productId := "1421094387"
	baseAddr := "https://www.ozon.ru/product/"
	addr := baseAddr + productId
	c.Visit(addr)
}
