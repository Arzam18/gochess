# GoChess Android ARM64 build notes

## Two separate Android builds

The workflow produces two independent executables and artifacts:

- **NEON:** `GOARM64=v8.0`, baseline ARM64 Advanced SIMD, intended for ordinary ARM64 Android devices.
- **DotProd:** `GOARM64=v8.2` plus the `gochess_dotprod` build tag and NDK `-march=armv8.2-a+dotprod`. This executable contains actual `SDOT` instructions and requires a CPU that supports the ARM DotProd extension.

Use the NEON executable on unknown or older ARM64 devices. The DotProd executable is suitable for the Samsung Galaxy A16 5G, whose ARM64-v8.2-A CPU supports DotProd, but should not be copied to a device without that feature.

## SIMD implementation details

The original `nnue_arm64.s` remains the baseline-NEON reference. In the explicit DotProd variant, two int8 NNUE matrix kernels use additional ARMv8.2-A DotProd instructions:

1. `nnueMatMul32x512_i8DotProd` uses `SDOT` directly. Its input comes from CReLU packing and is clamped to 0–127, so the input bytes are representable as signed int8 without changing their values.
2. `nnueV5L1Int8MatMulNDotProd` handles the SCReLU-packed unsigned 0–255 input. It maps each byte `x` to signed `x-128`, performs `SDOT`, then adds the exact correction `128 * sum(weights)`. This preserves the original unsigned-input × signed-weight dot product, including inputs above 127.

The ordinary NEON implementations are retained under `nnueMatMul32x512_i8NEON` and `nnueV5L1Int8MatMulNNEON`. The `gochess_dotprod` build wrappers choose the new SDOT kernels; the default build wrappers continue to call the NEON reference functions.

The remaining NNUE operations continue to use the existing NEON assembly. This is a targeted DotProd acceleration, not a wholesale AVX-512-equivalent ARM backend.

## Android build configuration

Both variants use:

- `GOOS=android`, `GOARCH=arm64`
- Android NDK r27c, API 24 Clang target
- `CGO_ENABLED=1`, retaining Fathom/Syzygy support
- Android PIE executable linked against `/system/bin/linker64`

The workflow runs host Go tests, cross-compiles each variant, verifies the AArch64 Android ELF interpreter, and checks the generated machine code. The DotProd job fails unless `SDOT` is found in the actual executable disassembly. It also packages a device-runnable differential test binary that compares the DotProd matrix kernels against their retained NEON reference implementations.

## Device verification

For the DotProd artifact, run the included test executable on a DotProd-capable Android device:

```sh
chmod +x gochess-dotprod-kernel-tests
./gochess-dotprod-kernel-tests -test.run 'TestDotProd'
```

Then run the engine through UCI and verify NNUE loading, search, NPS, temperature, and match strength. A successful compile and SDOT disassembly check do not establish playing strength. Keep GoChess in **Testing** until phone verification and MCEC matches are complete.
