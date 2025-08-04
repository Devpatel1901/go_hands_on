package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func hrefs(r io.Reader, base string) []string {
	links := parseHTML(r)

	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.href, "/"):
			ret = append(ret, base+l.href)
		case strings.HasPrefix(l.href, "http"):
			ret = append(ret, l.href)
		}
	}

	return ret
}

func filter(links []string, prefix string) []string {
	var ret []string

	for _, l := range links {
		if strings.HasPrefix(l, prefix) {
			ret = append(ret, l)
		}
	}

	return ret
}

func get(url string) []string {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Not able to reach the page pointed by url: %v", url)
		return []string{}
	}

	defer resp.Body.Close()

	reqUrl := resp.Request.URL

	baseURL := reqUrl.Scheme + "://" + reqUrl.Host

	return filter(hrefs(resp.Body, baseURL), baseURL)
}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})

	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})

		if len(q) == 0 {
			break
		}

		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}

			seen[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}
func main() {
	urlString := flag.String("url", "https://gophercises.com", "the url that you want to crawl")
	maxDepth := flag.Int("depth", 10, "the maximum number of links deep to traverse")

	flag.Parse()

	pages := bfs(*urlString, *maxDepth)

	toXml := urlset{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()
}
