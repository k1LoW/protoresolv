package protoresolv

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bufbuild/protocompile"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var _ protocompile.Resolver = (*Resolver)(nil)

type Resolver struct {
	importPaths []string
	sources     map[string][]byte
	mu          sync.RWMutex
}

type Option func(*Resolver) error

func Proto(protos ...string) Option {
	return func(r *Resolver) error {
		const sep = string(filepath.Separator)
		for _, p := range protos {
			abs, err := filepath.Abs(p)
			if err != nil {
				return err
			}
			resolved := false
			for _, ip := range r.importPaths {
				if strings.HasPrefix(abs, ip+sep) {
					path := strings.TrimPrefix(abs, ip+sep)
					if _, ok := r.sources[path]; ok {
						resolved = true
						break
					}
					b, err := os.ReadFile(p)
					if err != nil {
						return err
					}
					r.mu.Lock()
					r.sources[path] = b
					r.mu.Unlock()
					resolved = true
					break
				}
			}
			if !resolved {
				b, err := os.ReadFile(p)
				if err != nil {
					return err
				}
				r.mu.Lock()
				r.sources[p] = b
				r.mu.Unlock()
			}
		}
		return nil
	}
}

// New creates a new resolver.
func New(importPaths []string, opts ...Option) (*Resolver, error) {
	r := &Resolver{
		sources: map[string][]byte{},
	}
	var resolvedImportPaths []string
	for _, importPath := range importPaths {
		abs, err := filepath.Abs(importPath)
		if err != nil {
			return nil, err
		}
		if fi, ok := os.Stat(abs); ok != nil || !fi.IsDir() {
			return nil, fmt.Errorf("import path %q is not a directory", importPath)
		}
		resolvedImportPaths = append(resolvedImportPaths, abs)
	}

	for _, resolvedImportPath := range resolvedImportPaths {
		sources := map[string][]byte{}
		if err := filepath.WalkDir(resolvedImportPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".proto" {
				return nil
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(resolvedImportPath, path)
			if err != nil {
				return err
			}
			sources[rel] = b
			return nil
		}); err != nil {
			return nil, err
		}
		r.mu.Lock()
		for path, b := range sources {
			r.sources[path] = b
		}
		r.mu.Unlock()
	}

	r.importPaths = resolvedImportPaths
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Paths returns the paths of the file descriptor sets and sources.
func (r *Resolver) Paths() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var paths []string
	for path := range r.sources {
		paths = append(paths, path)
	}
	return paths
}

func (r *Resolver) FindFileByPath(path string) (protocompile.SearchResult, error) {
	result := protocompile.SearchResult{}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if b, ok := r.sources[path]; ok {
		result.Source = bytes.NewReader(b)
		return result, nil
	}
	return result, protoregistry.NotFound
}
