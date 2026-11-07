# Android ARM64 Build

The GitHub Actions workflow produces two standalone Android ARM64 UCI executables:

- `GoChess-Android-ARM64-NEON` — baseline ARM64 NEON (`GOARM64=v8.0`).
- `GoChess-Android-ARM64-DotProd` — ARMv8.2-A DotProd (`GOARM64=v8.2`), with real SDOT kernels.

## Embedded production NNUE

Before compiling either variant, the workflow downloads the exact production network referenced by `net.txt`:

`net-v5-120sb-sb120.nnue`

Source: `https://github.com/adamtwiss/gochess/releases/download/v0.5.0/net-v5-120sb-sb120.nnue`

The workflow checks the downloaded file's SHA-256, the production MD5 recorded in `models.md`, and minimum size, then compiles it into the executable with the `gochess_embed_net` build tag and Go's `embed` support. The NNUE is parsed directly from the embedded byte slice at startup; no temporary network file is extracted on Android.

The two GitHub Actions artifacts each contain **exactly one file**, `gochess`. The executable includes the production NNUE internally, so a separate `.nnue` or `net.txt` file is not required. The opening book is not embedded and remains optional.

External network files still take priority over the embedded production network when found beside the executable or in the current working directory. The UCI `NNUEFile` option continues to permit an explicit network override.

## Rebuilding locally

The source archive intentionally does not include the large release network. To create an embedded build manually:

1. Download the official network to `assets/net-v5-120sb-sb120.nnue`.
2. Verify its SHA-256 against `embeddedNNUEExpectedSHA256` in `nnue_embedded.go`.
3. Build with `-tags gochess_embed_net` (and add `gochess_dotprod` for the DotProd variant).

Without `gochess_embed_net`, the source retains its original external `net.txt` / `NNUEFile` network-loading behavior.

## Verification performed by Actions

- Download and SHA-256 verification of the production NNUE.
- Normal Go test suite and embedded-network test suite.
- Host UCI handshake from an isolated working directory, verifying that the embedded network loads without external files.
- Android AArch64 ELF and `/system/bin/linker64` checks.
- NEON instruction verification in both binaries.
- SDOT disassembly verification in the DotProd binary only.
- Artifact check that `dist` contains only the `gochess` executable.

Physical Android runtime, NPS and playing-strength verification are still required after downloading the artifacts.
