package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

type Links struct {
	data []LinkContent
}

func (l Links) Print() {
	for i, link := range l.data {
		fmt.Printf("Link %d:\n", i+1)
		fmt.Printf("  Href: %s\n", link.href)
		fmt.Printf("  Text: %s\n", link.text)
	}
}

type LinkContent struct {
	href string
	text string
}

func main() {
	htmlFileName := flag.String("htmlFile", "ex1.html", "html file path")

	flag.Parse()

	file, err := os.Open(*htmlFileName)

	if err != nil {
		log.Printf("Unable to read file %v", file)
		return
	}

	defer file.Close()

	z := html.NewTokenizer(file)

	depth := 0
	globalLinkContent := LinkContent{href: "", text: ""}
	links := Links{
		data: []LinkContent{},
	}

	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			break
		}

		switch tt {
		case html.StartTagToken:
			tn, hasAttr := z.TagName()

			if len(tn) == 1 && tn[0] == 'a' {
				depth++

				if depth == 1 {
					if hasAttr {
						_, value, _ := z.TagAttr()

						globalLinkContent.href = string(value)
					}
				}
			}
		case html.TextToken:
			token := z.Token()

			if depth == 1 {
				globalLinkContent.text += token.Data
			}
		case html.EndTagToken:
			tn, _ := z.TagName()

			if len(tn) == 1 && tn[0] == 'a' {
				if depth == 1 {
					links.data = append(links.data, globalLinkContent)
					globalLinkContent = LinkContent{href: "", text: ""}
				}
				depth--
			}
		case html.SelfClosingTagToken:
			tn, _ := z.TagName()

			if depth == 0 && len(tn) == 1 && tn[0] == 'a' {
				_, value, _ := z.TagAttr()

				globalLinkContent.href = string(value)

				links.data = append(links.data, globalLinkContent)
				globalLinkContent = LinkContent{href: "", text: ""}
			}
		}
	}

	links.Print()
}
