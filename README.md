# Local proto file resolver for github.com/bufbuild/protocompile [![Go Reference](https://pkg.go.dev/badge/github.com/k1LoW/protoresolv.svg)](https://pkg.go.dev/github.com/k1LoW/protoresolv) [![build](https://github.com/k1LoW/protoresolv/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/protoresolv/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/protoresolv/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/protoresolv/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/protoresolv/time.svg)

## Usage

``` go
importPaths := []string{"path/to/your/proto"}
r, _ := protoresolv.New(protoresolv.New(importPaths)
comp := protocompile.Compiler{
	Resolver: protocompile.WithStandardImports(r),
}
fds, _ := comp.Compile(ctx, r.Paths()...)
```
