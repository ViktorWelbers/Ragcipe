package db

import (
	"context"

	"github.com/qdrant/go-client/qdrant"
)

type Qdrant struct {
	client *qdrant.Client
}

func NewClient() (Qdrant, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		return Qdrant{}, err
	}

	return Qdrant{client: client}, nil
}

func (q *Qdrant) CreateCollection() {
	q.client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "Recipes",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     4,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}
