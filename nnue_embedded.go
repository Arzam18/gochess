//go:build gochess_embed_net

package chess

import _ "embed"

// The workflow downloads and verifies this asset before compiling the two
// Android engine variants. The binary is therefore self-contained at runtime.
//
//go:embed assets/net-v5-120sb-sb120.nnue
var embeddedNNUE []byte

const embeddedNNUEFilename = "net-v5-120sb-sb120.nnue"
const embeddedNNUEExpectedSHA256 = "be22a20bfde78d7915cc0cc5fe84e917a76d5df130b9555884a570a04b3ab2fc"
const embeddedNNUEExpectedMD5 = "4a58e81ac236f4f63e7ffecc7c0a2ddc"
