//go:build arm64 && gochess_dotprod

package chess

import (
	"math/rand"
	"testing"
)

func TestDotProdClassicInt8MatMulMatchesNEON(t *testing.T) {
	rng := rand.New(rand.NewSource(0x474f4348))
	input := make([]byte, 512)
	weights := make([]int8, 32*512)
	biases := make([]int32, 32)
	for i := range input {
		input[i] = byte(rng.Intn(128))
	}
	for i := range weights {
		weights[i] = int8(rng.Intn(256) - 128)
	}
	for i := range biases {
		biases[i] = int32(rng.Intn(200001) - 100000)
	}
	want, got := make([]int32, 32), make([]int32, 32)
	nnueMatMul32x512_i8NEON(&input[0], &weights[0], &biases[0], &want[0])
	nnueMatMul32x512_i8DotProd(&input[0], &weights[0], &biases[0], &got[0])
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("classic output %d: DotProd=%d NEON=%d", i, got[i], want[i])
		}
	}
}

func TestDotProdV5Int8MatMulMatchesNEON(t *testing.T) {
	rng := rand.New(rand.NewSource(0x56354938))
	const accLen, l1 = 768, 64
	input := make([]byte, accLen)
	weights := make([]int8, accLen*l1)
	for i := range input {
		input[i] = byte(rng.Intn(256))
	}
	for i := range weights {
		weights[i] = int8(rng.Intn(256) - 128)
	}
	want, got := make([]int32, l1), make([]int32, l1)
	for i := range want {
		want[i] = int32(rng.Intn(200001) - 100000)
		got[i] = want[i]
	}
	nnueV5L1Int8MatMulNNEON(&input[0], &weights[0], &want[0], accLen, l1)
	nnueV5L1Int8MatMulNDotProd(&input[0], &weights[0], &got[0], accLen, l1)
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("v5 output %d: DotProd=%d NEON=%d", i, got[i], want[i])
		}
	}
}
