package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func buildAvitoURL(page int) string {
	u, _ := url.Parse("https://www.avito.ru/sankt-peterburg/predlozheniya_uslug/oborudovanie_proizvodstvo-ASgBAgICAUSYC7SfAQ")
	q := u.Query()
	q.Set("p", strconv.Itoa(page))
	q.Set("q", "аренда генераторов")
	u.RawQuery = q.Encode()
	return u.String()
}

func ExampleScrape() {
	client := &http.Client{}

	for page := 1; page <= 10; page++ {
		link := buildAvitoURL(page)
		fmt.Println("PAGE:", page)

		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9")

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find(`[data-marker="item"]`).Each(func(i int, s *goquery.Selection) {
			title := s.Find(`[itemprop="name"]`).Text()
			price, ok := s.Find(`meta[itemprop="price"]`).Attr("content")
			if !ok {
				return
			}
			link, _ := s.Find("a").Attr("href")

			fmt.Println(title, price+"₽", "https://www.avito.ru"+link)
		})
	}
}

func main() {
	ExampleScrape()
}
