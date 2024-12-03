package main

import (
	"fmt"
	"gorecipe/pkg/db"
	"gorecipe/pkg/llm"
	"log"
)

func main() {
	qdrant, err := db.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	// qdrant.CreateCollection()
	embedding := llm.CreateEmbeddings("chicken burger")
	recipes, err := qdrant.QueryVector(embedding)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(recipes)
}
