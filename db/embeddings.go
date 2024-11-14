package db

import (
	"context"

	"github.com/qdrant/go-client/qdrant"
)

type Qdrant struct {
	client *qdrant.Client
}

func NewClient() (*Qdrant, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		return &Qdrant{}, err
	}

	return &Qdrant{client: client}, nil
}

func (q *Qdrant) CreateCollection() {
	q.client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "Recipes",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1024,
			Distance: qdrant.Distance_Cosine,
		}),
	})
}

func (q *Qdrant) InsertVector(embedding []float64, embeddingData map[string]any) (*qdrant.UpdateResult, error) {
	embeddingVectors := make([]float32, len(embedding))
	for i, v := range embedding {
		embeddingVectors[i] = float32(v)
	}
	operationInfo, err := q.client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "Recipes",
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(1),
				Vectors: qdrant.NewVectors(embeddingVectors...),
				Payload: qdrant.NewValueMap(embeddingData),
			},
		},
	})
	return operationInfo, err
}

func (q *Qdrant) QueryVector(embedding []float64) (*qdrant.UpdateResult, error) {
	embeddingVectors := make([]float32, len(embedding))
	for i, v := range embedding {
		embeddingVectors[i] = float32(v)
	}
	operationInfo, err := q.client.QueryBatch(context.Background(), &qdrant.QueryBatchPoints{
		CollectionName: "Recipes",
		QueryPoints: []*qdrant.QueryPoints{
			{
				Id:      qdrant.NewIDNum(1),
				Vectors: qdrant.NewVectors(embeddingVectors...),
				Payload: qdrant.NewValueMap(embeddingData),
			},
		},
	})
	return operationInfo, err
}
