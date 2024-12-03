package main

import (
	"fmt"
	"gorecipe/pkg/db"
	"gorecipe/pkg/llm"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	qdrant *db.Qdrant
}

func (h *Handler) recipeHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	embedding := llm.CreateEmbeddings(string(body))
	recipes, err := h.qdrant.QueryVector(embedding)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(recipes)
}

func main() {
	qdrant, err := db.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	handler := Handler{qdrant: qdrant}
	router := http.NewServeMux()
	router.HandleFunc(http.MethodPost+" /", handler.recipeHandler)
	fmt.Println("Starting server at port 8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Server failed:", err)
	}
}
