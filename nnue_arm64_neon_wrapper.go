//go:build arm64 && !gochess_dotprod

package chess

// Default ARM64 build: retain the baseline-NEON implementations.
func nnueMatMul32x512_i8(input *byte, weightsT *int8, biases *int32, output *int32) {
	nnueMatMul32x512_i8NEON(input, weightsT, biases, output)
}

func nnueV5L1Int8MatMulN(acc8 *byte, wT8 *int8, hidden *int32, accLen int, l1 int) {
	nnueV5L1Int8MatMulNNEON(acc8, wT8, hidden, accLen, l1)
}
