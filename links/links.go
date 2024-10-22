package links

import (
	"log"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Scraper struct {
	mu sync.Mutex
}

func (s *Scraper) processLink(n *html.Node) {
	for _, a := range n.Attr {
		if a.Key == "href" && strings.Contains(a.Val, "/rezepte/") {
			href := a.Val
			if !strings.Contains(href, "https://") {
				s.mu.Lock()
				defer s.mu.Unlock()
				href = "https://www.rewe.de" + href
				f, err := os.OpenFile("links.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
				_, err = f.WriteString(href)
				_, err = f.WriteString("\n")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func FetchAlLRecipeLinks(data string, wg *sync.WaitGroup) {
	var processAllProduct func(*html.Node)
	scraper := Scraper{}
	processAllProduct = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			scraper.processLink(n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processAllProduct(c)
		}
	}
	n, err := html.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	processAllProduct(n)
	wg.Done()
}
