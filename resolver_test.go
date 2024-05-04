package protoresolv

import (
	"slices"
	"testing"
)

func TestImportPaths(t *testing.T) {
	tests := []struct {
		importPaths []string
		want        []string
	}{
		{[]string{"testdata/helloapis"}, []string{"acme/hello/v2/hello.proto"}},
		{[]string{"testdata/helloapis/acme"}, []string{"hello/v2/hello.proto"}},
	}
	for _, tt := range tests {
		r, err := New(tt.importPaths)
		if err != nil {
			t.Fatal(err)
		}
		paths := r.Paths()
		for _, p := range tt.want {
			if !slices.Contains(paths, p) {
				t.Errorf("got %v, want %v", paths, p)
			}
		}
	}
}

func TestProto(t *testing.T) {
	tests := []struct {
		importPaths []string
		protos      []string
		want        []string
	}{
		{
			[]string{"testdata/helloapis"},
			[]string{"testdata/helloapis/acme/hello/v2/hello.proto"},
			[]string{"acme/hello/v2/hello.proto"},
		},
		{
			[]string{"testdata"},
			[]string{"testdata/helloapis/acme/hello/v2/hello.proto"},
			[]string{"helloapis/acme/hello/v2/hello.proto"},
		},
		{
			[]string{""},
			[]string{"testdata/helloapis/acme/hello/v2/hello.proto"},
			[]string{"testdata/helloapis/acme/hello/v2/hello.proto"},
		},
	}
	for _, tt := range tests {
		r, err := New(tt.importPaths, Proto(tt.protos...))
		if err != nil {
			t.Fatal(err)
		}
		paths := r.Paths()
		for _, p := range tt.want {
			if !slices.Contains(paths, p) {
				t.Errorf("got %v, want %v", paths, p)
			}
		}
	}
}

func TestPaths(t *testing.T) {
	r, err := New([]string{"testdata/helloapis"})
	if err != nil {
		t.Fatal(err)
	}
	got := len(r.Paths())
	r2, err := New([]string{"testdata/helloapis"})
	if err != nil {
		t.Fatal(err)
	}
	got2 := len(r2.Paths())
	if got != got2 {
		t.Errorf("got %d, want %d", got2, got)
	}
}
