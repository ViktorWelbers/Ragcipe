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
	recipe := Recipe{}
	err = recipe.fromHTML(n)
	if err != nil {
		log.Fatal(err)
	}
	wg.Done()
	fmt.Println(recipe)
}

type Recipe struct {
	servings    string
	ingredients []Ingredient
}

type Ingredient struct {
	item   string
	amount int
}

func (recipe *Recipe) fromHTML(n *html.Node) error {
	if n.Data == "span" {
		for _, a := range n.Attr {
			text := a.Val
			if strings.Contains(text, "ld-rds text-base") {
				proposedServings := n.FirstChild.Data
				if n.NextSibling.Type == html.ElementNode && len(n.NextSibling.Data) > 1 {
					for _, b := range n.NextSibling.Attr {
						if b.Val == "getFormattedServingType()" {
							recipe.servings = proposedServings
						}
						// TODO: Get the servings type
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
					amount, err := strconv.Atoi(amountStr)
					if err != nil {
						return err
					}
					if c.NextSibling.Type == html.TextNode && len(c.NextSibling.Data) > 1 {
						item := c.NextSibling.Data
						ingredient := Ingredient{item, amount}
						recipe.ingredients = append(recipe.ingredients, ingredient)
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := recipe.fromHTML(c)
		if err != nil {
			return err
		}
	}

	return nil
}
