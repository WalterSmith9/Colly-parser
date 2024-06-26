package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly" //scraping framework
)

func main() {

	c := colly.NewCollector(colly.AllowedDomains("www.ozon.ru"))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Link of the page:", r.URL)
	})

	//getting description
	c.OnHTML("div.m7m_27", func(h *colly.HTMLElement) {
		name := h.ChildText("a.mm8_27")
		fmt.Println("Name: ", name)
	})

	//getting bread crumbs (path)
	c.OnHTML("ol.eg_10", func(h *colly.HTMLElement) {
		var breadCrumbs string
		h.ForEach(".ah6.ge0_10", func(i int, h *colly.HTMLElement) {
			breadCrumbs += h.ChildText("span")
			breadCrumbs += "--"
		})
		breadCrumbs = strings.Trim(breadCrumbs, "-")
		fmt.Println(breadCrumbs)
	})

	//getting characteristics
	c.OnHTML("div.d3.c5", func(h *colly.HTMLElement) {
		h.ForEach(".p0k_27", func(i int, h *colly.HTMLElement) {
			key := h.ChildText("dt.pk_27")
			val := h.ChildText("dd.kp0_27")
			fmt.Printf("%s - %s\n", key, val)
		})

	})

	//testing: writing a DOM into .html file
	// file, err := os.Create("pages/phone-char.html")
	// if err != nil {
	// 	fmt.Println("Unable to create file:", err)
	// 	os.Exit(1)
	// }
	// defer file.Close()

	// c.OnHTML("html", func(h *colly.HTMLElement) {
	// 	h.DOM.Each(func(i int, s *goquery.Selection) {
	// 		e, err := s.Html()
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		file.WriteString("<html>")
	// 		file.WriteString(e)
	// 		file.WriteString("</html>")
	// 		//fmt.Println(e)
	// 	})

	// })

	//setting url for parsing
	productId := "1467382973"
	addr := fmt.Sprintf("https://www.ozon.ru/product/%s/features/", productId)
	//addr := baseAddr + productId

	c.Visit(addr)
}
