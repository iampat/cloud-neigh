package vector

import (
	"log"
	"math/rand"
	"time"
)

type Vector32 struct {
	Values []float32 `json:"values"`
}

func (v *Vector32) Len() int {
	return len(v.Values)
}

func Dot(v1 *Vector32, v2 *Vector32) float32 {
	if v1.Len() != v2.Len() {
		log.Panicf("v1 and v2 mush have the same size len(v1)=%d len(v2)=%d\n", v1.Len(), v2.Len())
	}
	var dot float32 = 0.0
	for idx, val := range v1.Values {
		dot += v2.Values[idx] * val
	}
	return dot
}

func NewRandomVec(dim int, rnd *rand.Rand) *Vector32 {
	if rnd == nil {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	vec := Vector32{}
	for idx := 0; idx < dim; idx++ {
		vec.Values = append(vec.Values, rnd.Float32())
	}
	return &vec
}
