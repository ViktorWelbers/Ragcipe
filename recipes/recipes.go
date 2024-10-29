package recipes

import (
	"context"
	"fmt"
	"gorecipe/db"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/net/html"
)

type Ingredient struct {
	Name   string
	Amount int
}

type Recipe struct {
	servings     string
	servingsType string
	ingredients  []Ingredient
	instructions []string
}

func FetchRecipe(recipeUrl string, data string, wg *sync.WaitGroup, queries *db.Queries) {
	n, err := html.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	recipe := Recipe{}
	err = recipe.parseRecipe(n)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(recipe)
	servingsInt, err := strconv.Atoi(recipe.servings)
	if err != nil {
		log.Println("There was an error trying to convert servings:", recipe.servings)
	}
	rawUrl, err := url.Parse(recipeUrl)
	if err != nil {
		return
	}
	host := rawUrl.Hostname()
	splitHost := strings.Split(host, ".")
	countryCode := splitHost[len(splitHost)-1]
	createRecipeParams := db.CreateRecipeParams{
		Title:        "test",
		Servings:     int32(servingsInt),
		ServingsType: recipe.servingsType,
		CountryCode:  countryCode,
		HostUrl:      host,
		OriginalUrl: pgtype.Text{
			String: recipeUrl,
			Valid:  true,
		},
	}
	queries.CreateRecipe(context.Background(), createRecipeParams)
	wg.Done()
}

func (recipe *Recipe) extractInstructions(n *html.Node) error {
	if n.Type == html.ElementNode && n.Data == "div" {
		// Check if it has the text formatting classes that indicate an instruction div
		isInstruction := false
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "ld-rds mt- self-stretch text-sm leading-5 text-gray-1000 md:text-base lg:mt-6 lg:leading-6") {
				isInstruction = true
				break
			}
		}

		if isInstruction {
			// Find and extract the paragraph text
			paragraphNode := findFirstChild(n, "p")
			if paragraphNode != nil {
				instructionText := getTextContent(paragraphNode)
				if instructionText != "" {
					recipe.instructions = append(recipe.instructions, instructionText)
				}
			}
		}
	}
	return nil
}

func (recipe *Recipe) extractIngredients(n *html.Node) error {
	if n.Type == html.ElementNode && n.Data == "li" {
		// Check if it has the ingredient_list_item class
		isIngredient := false
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "ingredient_list_item") {
				isIngredient = true
				break
			}
		}

		if isIngredient {
			var ingredient Ingredient

			// Find the div containing the ingredient info
			divNode := findFirstChild(n, "div")
			if divNode == nil {
				return nil
			}

			// Get all text content
			fullText := getTextContent(divNode)
			fullText = strings.TrimSpace(fullText)

			// If there's a span with amount, extract it
			spanNode := findFirstChild(divNode, "span")
			if spanNode != nil && spanNode.FirstChild != nil {
				amountStr := spanNode.FirstChild.Data
				amount, err := strconv.Atoi(amountStr)
				if err == nil {
					ingredient.Amount = amount
					// Remove the amount from full text
					fullText = strings.TrimPrefix(fullText, amountStr)
				}
			}

			// Everything else is the ingredient name
			ingredient.Name = strings.TrimSpace(fullText)

			if ingredient.Name != "" {
				recipe.ingredients = append(recipe.ingredients, ingredient)
			}
		}
	}
	return nil
}

func (recipe *Recipe) extractServings(n *html.Node) error {
	if n.Type == html.ElementNode && n.Data == "div" {
		// Check if it has the interactive-element class
		isServingsDiv := false
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "interactive-element") {
				isServingsDiv = true
				break
			}
		}

		if isServingsDiv {
			// Find the span containing the servings info
			middleSpan := findFirstChild(n, "span")
			if middleSpan != nil {
				// Find the servings amount span
				servingsSpan := findFirstChild(middleSpan, "span")
				if servingsSpan != nil {
					recipe.servings = getTextContent(servingsSpan)

					// Find the servings type in the next sibling
					if servingsSpan.NextSibling != nil {
						recipe.servingsType = getTextContent(servingsSpan.NextSibling)
					}
				}
			}
		}
	}
	return nil
}

func findFirstChild(n *html.Node, tag string) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			return c
		}
	}
	return nil
}

func getTextContent(n *html.Node) string {
	var result string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			result += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return strings.TrimSpace(result)
}

func (recipe *Recipe) parseRecipe(n *html.Node) error {
	// Try to extract servings first
	if err := recipe.extractServings(n); err != nil {
		return err
	}

	// Try to extract ingredients
	if err := recipe.extractIngredients(n); err != nil {
		return err
	}
	if err := recipe.extractInstructions(n); err != nil {
		return err
	}
	// Continue traversing
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := recipe.parseRecipe(c)
		if err != nil {
			return err
		}
	}

	return nil
}
