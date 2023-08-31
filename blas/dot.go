package blas

func DotProduct(x, y []float32) float32 {
	var z float32 = 0.0
	for i, xi := range x {
		z += xi * y[i]
	}
	return z
}
