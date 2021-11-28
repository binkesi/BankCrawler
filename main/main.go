package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func JianhangSearch() {
	client := &http.Client{}
	url := "http://www1.ccb.com/cn/home/map/branchSearch.html"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Host", "www1.ccb.com")
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("Referer", "http://www1.ccb.com/cn/home/map/branchSearch.html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	fmt.Println(doc.Text())
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".city_list2 j_atm_list .li_a").Each(func(i int, s *goquery.Selection) {
		content := s.Find(".click_show").Text()
		fmt.Printf("%d: %s\n", i, content)
	})
}

func main() {
	JianhangSearch()
}
