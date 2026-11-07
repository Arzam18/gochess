//go:build gochess_embed_net

package chess

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"testing"
)

func TestEmbeddedProductionNNUE(t *testing.T) {
	if len(embeddedNNUE) < 20*1024*1024 {
		t.Fatalf("embedded production NNUE unexpectedly small: %d bytes", len(embeddedNNUE))
	}
	hash := sha256.Sum256(embeddedNNUE)
	gotHash := hex.EncodeToString(hash[:])
	if gotHash != embeddedNNUEExpectedSHA256 {
		t.Fatalf("embedded NNUE SHA-256 mismatch: got %s, want %s", gotHash, embeddedNNUEExpectedSHA256)
	}
	md5Hash := md5.Sum(embeddedNNUE)
	gotMD5 := hex.EncodeToString(md5Hash[:])
	if gotMD5 != embeddedNNUEExpectedMD5 {
		t.Fatalf("embedded NNUE MD5 mismatch: got %s, want %s", gotMD5, embeddedNNUEExpectedMD5)
	}
	if binary.LittleEndian.Uint32(embeddedNNUE[:4]) != nnueMagic {
		t.Fatal("embedded NNUE has invalid magic")
	}
	version := binary.LittleEndian.Uint32(embeddedNNUE[4:8])
	if version != nnueVersionV5 {
		t.Fatalf("unexpected embedded production NNUE version %d, want v5", version)
	}
	net, err := LoadNNUEV5FromBytes(embeddedNNUE)
	if err != nil {
		t.Fatalf("load embedded production NNUE: %v", err)
	}
	if net.HiddenSize != 1024 || net.UseSCReLU || net.UsePairwise || net.L1Size != 0 {
		t.Fatalf("unexpected production model architecture: hidden=%d SCReLU=%t pairwise=%t L1=%d", net.HiddenSize, net.UseSCReLU, net.UsePairwise, net.L1Size)
	}
	if len(net.InputWeights) != NNUEInputSize*1024 || len(net.InputBiases) != 1024 {
		t.Fatalf("production NNUE input tensors have unexpected sizes: weights=%d biases=%d", len(net.InputWeights), len(net.InputBiases))
	}
	if len(net.OutputWeights) != NNUEOutputBuckets*2*1024 {
		t.Fatalf("production NNUE output tensor has unexpected size: %d", len(net.OutputWeights))
	}
}
