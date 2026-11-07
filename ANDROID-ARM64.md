# GoChess Android ARM64 build notes

## Engine-only artifacts

The workflow creates two separate downloadable artifacts. Each artifact contains exactly one file: the `gochess` Android executable.

- **NEON:** `GOARM64=v8.0`, baseline ARM64 Advanced SIMD.
- **DotProd:** `GOARM64=v8.2`, the `gochess_dotprod` build tag, and NDK `-march=armv8.2-a+dotprod`. The executable is verified to contain actual `SDOT` instructions and requires a CPU that supports ARM DotProd.

No NNUE network, `net.txt`, opening book, test executable, README, or build report is included in either download. This keeps the engine artifact to the executable only.

**NNUE note:** GoChess loads its NNUE network separately. With only the executable present, it can start the UCI protocol but will run without NNUE until a compatible network is supplied through the UCI `NNUEFile` option (or placed beside the executable together with the matching `net.txt`). The network is not embedded in the binary. An opening book is likewise optional and is not bundled.

Use the NEON executable on unknown or older ARM64 devices. The DotProd executable is appropriate for the Samsung Galaxy A16 5G, but should not be copied to a device without DotProd support.

## SIMD implementation details

The original `nnue_arm64.s` remains the baseline-NEON implementation. In the DotProd variant, two int8 NNUE matrix kernels use ARMv8.2-A DotProd instructions:

1. `nnueMatMul32x512_i8DotProd` uses `SDOT` for CReLU inputs clamped to 0–127.
2. `nnueV5L1Int8MatMulNDotProd` handles unsigned 0–255 SCReLU inputs using exact centering and correction before/after `SDOT`.

The ordinary NEON implementations remain available as reference implementations. The default build uses NEON; the `gochess_dotprod` build tag selects the SDOT kernels. This is targeted DotProd acceleration, not a wholesale replacement of the entire NNUE backend.

## Android build and verification

Both variants use `GOOS=android`, `GOARCH=arm64`, Android NDK r27c, API 24, and CGO enabled for Syzygy/Fathom support. The workflow runs host tests and a UCI protocol smoke test, cross-compiles both variants, verifies the AArch64 ELF interpreter `/system/bin/linker64`, and inspects the actual machine code. The DotProd job fails if no `SDOT` instruction is found; the baseline NEON job fails if `SDOT` is present.

A device-runnable DotProd differential test binary is generated in the Actions runner's temporary directory for validation, but it is deliberately excluded from the downloadable engine artifact. The test executable is not uploaded.

A successful build and instruction check do not establish Android runtime correctness or playing strength. Test UCI, external NNUE loading, NPS, temperature, and MCEC results on the phone before marking GoChess successful.
