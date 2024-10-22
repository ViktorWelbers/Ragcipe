package recipes

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func FetchRecipe(data string, wg *sync.WaitGroup) {
	n, err := html.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	err = processRecipe(n)
	if err != nil {
		log.Fatal(err)
	}
	wg.Done()
}

type Recipe struct {
	ingredients []Ingredient
	servings    int
}

type Ingredient struct {
	item   string
	amount int
}

func processRecipe(n *html.Node) error {
	var ingredients []Ingredient
	var servings int
	if n.Data == "span" {
		for _, a := range n.Attr {
			text := a.Val
			if strings.Contains(text, "ld-rds text-base") {
				proposedServings := n.FirstChild.Data
				if n.NextSibling.Type == html.ElementNode && len(n.NextSibling.Data) > 1 {
					for _, b := range n.NextSibling.Attr {
						if b.Val == "getFormattedServingType()" {
							servings, err := strconv.Atoi(proposedServings)
							if err != nil {
								return err
							}
							fmt.Println(servings)
						}
					}
				}
			}
		}
	}
	if n.Type == html.ElementNode && n.Data == "li" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "span" {
				if len(c.Attr) > 0 {
					amountStr := strings.Trim(strings.Trim(c.Attr[0].Val, "adjustedAmount("), ")")
					fmt.Println(amountStr)
					amount, err := strconv.Atoi(amountStr)
					if err != nil {
						return err
					}
					if c.NextSibling.Type == html.TextNode && len(c.NextSibling.Data) > 1 {
						item := c.NextSibling.Data
						fmt.Println(amount, item)
						ingredient := Ingredient{item, amount}
						ingredients = append(ingredients, ingredient)
					}
				}
			}
		}
	}
	if ingredients != nil && servings != 0 {
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processRecipe(c)
	}
	return nil
}
