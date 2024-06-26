package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly" //scraping framework
)

type product struct {
	id                 string
	name               string
	breadCrumbs        string
	webCharacteristics map[string]string
}

func (p product) print() {
	fmt.Printf("%s;\"%s\";\"%s\"\n", p.id, "Название товара", p.name)
	fmt.Printf("%s;\"%s\";\"%s\"\n", p.id, "Хлебные крошки", p.breadCrumbs)
	for k, v := range p.webCharacteristics {
		fmt.Printf("%s;\"%s\";\"%s\"\n", p.id, k, v)
	}
}

func (p product) Write(src string) {
	file, err := os.Create(src)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()

	outputText := fmt.Sprintf("%s;\"%s\";\"%s\"\n", p.id, "Название товара", p.name)
	outputText += fmt.Sprintf("%s;\"%s\";\"%s\"\n", p.id, "Хлебные крошки", p.breadCrumbs)
	for k, v := range p.webCharacteristics {
		outputText += fmt.Sprintf("%s;\"%s\";\"%s\"\n", p.id, k, v)
	}

	file.WriteString(outputText)

}

func parseItem(id string) product {
	//creating entity for specific item
	item := &product{}
	item.id = id
	item.webCharacteristics = make(map[string]string)

	//creating parsing entity for specific item
	c := colly.NewCollector(colly.AllowedDomains("www.ozon.ru"))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Link of the page:", r.URL)
		//fmt.Println(r.Headers)
	})

	//getting description
	c.OnHTML("div.m7m_27", func(h *colly.HTMLElement) {
		item.name = h.ChildText("a.mm8_27")
	})

	//getting bread crumbs (path)
	c.OnHTML("ol.eg_10", func(h *colly.HTMLElement) {
		var breadCrumbs string
		h.ForEach(".ah6.ge0_10", func(i int, h *colly.HTMLElement) {
			breadCrumbs += h.ChildText("span")
			breadCrumbs += "--"
		})
		item.breadCrumbs = strings.Trim(breadCrumbs, "-")
	})

	//getting characteristics
	c.OnHTML("div.d3.c5", func(h *colly.HTMLElement) {
		h.ForEach(".p0k_27", func(i int, h *colly.HTMLElement) {
			key := h.ChildText("dt.pk_27")
			val := h.ChildText("dd.kp0_27")
			item.webCharacteristics[key] = val
		})

	})

	//setting url for parsing
	addr := fmt.Sprintf("https://www.ozon.ru/product/%s/features/", id)

	err := c.Visit(addr)
	if err != nil {
		log.Println(err)
	}

	// item.print()

	return *item
}

func main() {

	productId := "1421094387"
	parseItem(productId).Write("item-lists/output.txt")

	// err := c.SetProxy("http://188.235.0.207:8181")
	// if err != nil {
	// 	log.Println(err)
	// }

	//testing: writing a DOM into .html file
	// file, err := os.Create("pages/phone-proxy.html")
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

}
