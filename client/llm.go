package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

func GenerateRecipe(data string) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	messages := []api.Message{
		{
			Role:    "system",
			Content: "Provide very brief, concise responses. All output should be JSON. In the format {recipe : RECIPE_NAME , ingredients: []}.",
		},
		{
			Role:    "user",
			Content: "Please make a shopping plan for a meal. I would love to eat meal that are high in protein.",
		},
		{
			Role:    "assistant",
			Content: "{ \"recipe\": \"chicken curry\", \"ingredients\": [\"chicken\", \"curry powder\", \"onion\", \"garlic\", \"ginger\", \"tomato\", \"coconut milk\"]}",
		},
		{
			Role:    "user",
			Content: "I would like a meal for another day. Please provide me with a recipe that uses lentils.",
		},
	}
	stream := false
	req := &api.ChatRequest{
		Stream:   &stream,
		Model:    "llama3.1:8b",
		Messages: messages,
		Format:   "json",
	}
	var response string

	respFunc := func(resp api.ChatResponse) error {
		response = resp.Message.Content
		return nil
	}

	ctx := context.Background()
	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}

func CreateEmbeddings(text string) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	req := &api.EmbeddingRequest{
		Model:  "mxbai-embed-large",
		Prompt: text,
	}

	resp, err := client.Embeddings(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Embedding)
}
