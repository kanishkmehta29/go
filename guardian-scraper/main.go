package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

func randomUserAgent() string {
	ua := userAgents[rand.Intn(6)]
	return ua
}

func getRequest(targetURL string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", targetURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())
	res, err := client.Do(req)
	return res, err
}

func checkRelative(href string, baseURL string) string {
	if strings.HasPrefix(href, "/") {
		return fmt.Sprintf("%s%s", baseURL, href)
	} else {
		return href
	}
}

func resolveRelativeLinks(href string, baseURL string) (bool, string) {
	resulthref := checkRelative(href, baseURL)
	baseParse, _ := url.Parse(baseURL)
	resultParse, _ := url.Parse(resulthref)
	if baseParse != nil && resultParse != nil {
		if baseParse.Host == resultParse.Host {
			return true, resulthref
		}
	}
	return false, ""
}

func discoverLinks(res *http.Response, baseURL string) []string {
	if res != nil {
		doc, _ := goquery.NewDocumentFromReader(res.Body)
		foundURLs := []string{}
		if doc != nil {
			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				new_url, _ := s.Attr("href")
				foundURLs = append(foundURLs, new_url)
			})
		}
		return foundURLs
	}
	return nil
}

var tokens = make(chan struct{}, 5)

func Crawl(targetURL string, baseURL string) []string {
	fmt.Println(targetURL)
	tokens <- struct{}{}
	resp, _ := getRequest(targetURL)
	<-tokens
	links := discoverLinks(resp, baseURL)
	foundUrls := []string{}

	for _, link := range links {
		ok, correctLink := resolveRelativeLinks(link, baseURL)
		if ok && correctLink != "" {
			foundUrls = append(foundUrls, correctLink)
		}
	}
	return foundUrls
}

func main() {
	baseDomain := "https://www.theguardian.com"
	worklist := make(chan []string)
	go func() {
		worklist <- []string{"https://www.theguardian.com"}
	}()

	seen := make(map[string]bool)
	n := 1
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string, baseURL string) {
					foundLinks := Crawl(link, baseDomain)
					if foundLinks != nil {
						worklist <- foundLinks
					}
				}(link, baseDomain)
			}
		}
	}

}
