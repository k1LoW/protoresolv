# Local proto file resolver for github.com/bufbuild/protocompile

## Usage

``` go
importPaths := []string{"path/to/your/proto"}
r, _ := protoresolv.New(protoresolv.New(importPaths)
comp := protocompile.Compiler{
	Resolver: protocompile.WithStandardImports(r),
}
fds, _ := comp.Compile(ctx, r.Paths()...)
```
