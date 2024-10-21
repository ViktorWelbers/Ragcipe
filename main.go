package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"golang.org/x/net/html"
)

func main() {
	if _, err := os.Stat("links.txt"); err == nil {
		f, err := os.Open("links.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			link := scanner.Text()
			scrapeUrl(link, fetchRecipe)
			break
		}
	} else if errors.Is(err, os.ErrNotExist) {
		sum := 1
		for {
			scrapeUrl(fmt.Sprintf("https://www.rewe.de/rezepte/?pageNumber=%d", sum), fetchAlLRecipeLinks)
			res, err := http.Get(fmt.Sprintf("https://www.rewe.de/rezepte/?pageNumber=%d", sum+1))
			if err != nil || res.StatusCode != 200 {
				break
			}
			sum++
		}
	} else {
	}
}

func scrapeUrl(url string, dataextractorFunc func(string)) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			data := string(r.Body)
			dataextractorFunc(data)
		},
	}).Start()
}

func processRecipe(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "li" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "span" {
				if len(c.Attr) > 0 {
					amount := c.Attr[0].Val
					if c.NextSibling.Type == html.TextNode && len(c.NextSibling.Data) > 1 {
						ingredient := c.NextSibling.Data
						fmt.Println(amount, ingredient)
					}
				}
			}
		}
	}
	// TODO: EXTRACT SERVING SIZE
	if n.Type == html.ElementNode && n.Data == "span" {
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processRecipe(c)
	}
}

func fetchRecipe(data string) {
	n, err := html.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	processRecipe(n)
}

func processLink(n *html.Node) {
	for _, a := range n.Attr {
		if a.Key == "href" && strings.Contains(a.Val, "/rezepte/") {
			href := a.Val
			if !strings.Contains(href, "https://") {
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

func fetchAlLRecipeLinks(data string) {
	var processAllProduct func(*html.Node)
	processAllProduct = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			processLink(n)
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
}
