package openai

import "github/iampat/cloudy-neigh/vector"

type Embedder interface {
	Embeddings([]string) ([]*vector.Vector32, error)
	EmbeddingsWithCost([]string) ([]*vector.Vector32, int, error)
	EmbeddingDim() int
}
