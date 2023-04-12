package embeddings

import "github/iampat/cloudy-neigh/vector"

type Embedding interface {
	Get([]string) ([]vector.Vector32, error)
}
