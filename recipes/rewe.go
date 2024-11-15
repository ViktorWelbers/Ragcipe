package recipes

import (
	"encoding/json"
	"fmt"
	"gorecipe/db"
	"gorecipe/llm"
	"log"
	"strings"

	"golang.org/x/net/html"
)

type Ingredient struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
}

type Recipe struct {
	Servings     string       `json:"servings"`
	ServingsType string       `json:"servingsType"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []string     `json:"instructions"`
}

func FetchRecipe(recipeUrl string, data string, queries *db.Qdrant) {
	n, err := html.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	recipe := Recipe{}
	err = recipe.parseRecipe(n)
	if err != nil {
		log.Fatal(err)
	}
	jsonRecipe, err := json.Marshal(recipe)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	embeddingData := map[string]interface{}{
		"recipe": string(jsonRecipe),
	}
	embeddings := llm.CreateEmbeddings(string(jsonRecipe))
	queries.InsertVector(embeddings, embeddingData)
	fmt.Println(recipe)
}

func (recipe *Recipe) extractInstructions(n *html.Node) error {
	if n.Type == html.ElementNode && n.Data == "div" {
		// Check if it has the text formatting classes that indicate an instruction div
		isInstruction := false
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "ld-rds") {
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
					recipe.Instructions = append(recipe.Instructions, instructionText)
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
				ingredient.Amount = amountStr
				// Remove the amount from full text
				fullText = strings.TrimPrefix(fullText, amountStr)
			}

			// Everything else is the ingredient name
			ingredient.Name = strings.TrimSpace(fullText)

			if ingredient.Name != "" {
				recipe.Ingredients = append(recipe.Ingredients, ingredient)
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
					recipe.Servings = getTextContent(servingsSpan)

					// Find the servings type in the next sibling
					if servingsSpan.NextSibling != nil {
						recipe.ServingsType = getTextContent(servingsSpan.NextSibling)
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
