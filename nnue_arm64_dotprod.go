//go:build arm64 && gochess_dotprod

package chess

// These kernels use ARMv8.2-A DotProd SDOT instructions. Build this variant
// only for devices whose CPUs advertise the DotProd feature.
//
//go:noescape
func nnueMatMul32x512_i8DotProd(input *byte, weightsT *int8, biases *int32, output *int32)

//go:noescape
func nnueV5L1Int8MatMulNDotProd(acc8 *byte, wT8 *int8, hidden *int32, accLen int, l1 int)

func nnueMatMul32x512_i8(input *byte, weightsT *int8, biases *int32, output *int32) {
	nnueMatMul32x512_i8DotProd(input, weightsT, biases, output)
}

func nnueV5L1Int8MatMulN(acc8 *byte, wT8 *int8, hidden *int32, accLen int, l1 int) {
	nnueV5L1Int8MatMulNDotProd(acc8, wT8, hidden, accLen, l1)
}
