package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly" //scraping framework
)

type product struct {
	id                 string
	name               string
	breadCrumbs        string
	webCharacteristics map[string]string
}

func (p product) Print() {
	fmt.Printf("%s;\"%s\";\"%s\"\n", p.id, "Название товара", p.name)
	fmt.Printf("%s;\"%s\";\"%s\"\n", p.id, "Хлебные крошки", p.breadCrumbs)
	for k, v := range p.webCharacteristics {
		fmt.Printf("%s;\"%s\";\"%s\"\n", p.id, k, v)
	}
}

func (p product) Write(src string) {
	file, err := os.OpenFile(src, os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("Unable to open file:", err)
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

func (p product) NotFoundWrite(src string) {
	file, err := os.OpenFile(src, os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("Unable to open file:", err)
		os.Exit(1)
	}
	defer file.Close()

	outputText := p.id + "\n"

	file.WriteString(outputText)
}

func parseItem(id string, tokens chan struct{}) (product, error) {
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

	tokens <- struct{}{}
	err := c.Visit(addr)
	<-tokens
	if err != nil {
		return *item, err

	}

	//check if exists
	if item.name == "" {
		err := errors.New("no name")
		return *item, err
	}

	// item.print()

	return *item, nil
}

func getIdList(src string) []string {
	list := make([]string, 0, 10)

	file, err := os.Open(src)
	if err != nil {
		fmt.Println("Unable to open file:", err)
		os.Exit(1)
	}
	defer file.Close()
	//reading from file

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	return list
}

func main() {

	if err := os.Truncate("item-lists/not_found.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	if err := os.Truncate("item-lists/output.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	productIdList := getIdList("item-lists/input.txt")

	wg := new(sync.WaitGroup)
	tokens := make(chan struct{}, 100)

	muNotFound := new(sync.Mutex)
	muOutput := new(sync.Mutex)

	for _, v := range productIdList {
		wg.Add(1)

		go func(v string) {
			item, err := parseItem(v, tokens)
			if err != nil {

				log.Println(item.id, err)
				muNotFound.Lock()
				item.NotFoundWrite("item-lists/not_found.txt")
				muNotFound.Unlock()

			} else {

				muOutput.Lock()
				item.Write("item-lists/output.txt")
				muOutput.Unlock()

			}
			wg.Done()

		}(v)
	}

	wg.Wait()

}
