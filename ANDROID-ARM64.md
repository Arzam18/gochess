# GoChess Android ARM64 build notes

## ARM/NEON audit

GoChess already contains a genuine ARM64 NEON NNUE backend:

- `nnue_arm64.go` selects the ARM64 SIMD implementation (`nnueUseSIMD = true` and `nnueUseSIMDV5 = true`).
- `nnue_arm64.s` contains the NEON accumulator, activation, and integer matrix/vector kernels.
- The ARM64 source uses baseline Advanced SIMD instructions such as `SMAX`, `SMIN`, `SMULL`, `SMLAL`, `SADALP`, and `SQXTUN`.
- No `SDOT`/`UDOT`, SVE, or other optional ARMv8.2+ vector instruction is required by this backend. ARM64 Android devices have the baseline Advanced SIMD/NEON support needed by these kernels.

The existing NEON kernels are retained. This build does not raise the whole binary's minimum CPU target to `armv8.2-a+dotprod`.

## Android-specific correction

The ARM64 NNUE code enabled `nnueUseSIMDV5`, but `nnueFloatMatVecFMA` was still implemented as a deliberate-crash assembly stub. The v5/v7 SCReLU float L2 path can call that function. The ARM64 version now implements the same scalar float32 matrix-vector loop as the generic Go fallback. This affects only that small float32 layer; the integer NNUE operations continue to use the existing NEON assembly.

## Build target

The workflow builds:

- `GOOS=android`
- `GOARCH=arm64`
- `GOARM64=v8.0`
- Android NDK r27c, API 24 Clang target
- `CGO_ENABLED=1`, preserving the bundled Fathom/Syzygy tablebase support

The workflow runs the host Go tests, cross-compiles the Android executable, checks its AArch64 ELF header, verifies that NEON instructions appear in disassembly, and packages the executable with `net.txt`, the NNUE network referenced by `net.txt`, `book.bin`, and this readme.

## Validation status

A successful GitHub Actions build verifies compilation and the presence of NEON instructions. It does not by itself verify UCI operation, NNUE correctness on a physical phone, NPS, heat behavior, or playing strength. Keep GoChess in **Testing** until the Android runtime and MCEC match checks are complete.
