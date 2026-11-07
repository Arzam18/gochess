//go:build !gochess_embed_net

package chess

// Non-embedded builds retain the original external net.txt / NNUEFile flow.
var embeddedNNUE []byte

const embeddedNNUEFilename = "net-v5-120sb-sb120.nnue"
const embeddedNNUEExpectedSHA256 = "be22a20bfde78d7915cc0cc5fe84e917a76d5df130b9555884a570a04b3ab2fc"
const embeddedNNUEExpectedMD5 = "4a58e81ac236f4f63e7ffecc7c0a2ddc"
