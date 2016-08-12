package flake

import (
	"sort"
	"testing"
)

func CreateNewFlake(t *testing.T) Flaker {
	flake, err := New()

	if err != nil {
		t.Errorf("Unable to create new Flake: %s", err)
	}
	return flake
}

func TestNewFlake(t *testing.T) {
	f := CreateNewFlake(t)

	var ids []string

	for i := 0; i < 4; i++ {
		id := f.NextId()

		ids = append(ids, id.String())
	}

	if !sort.StringsAreSorted(ids) {
		t.Errorf("IDs are not sorted!")
	}
}

func BenchmarkNextId(b *testing.B) {

	f, err := New()

	if err != nil {
		b.Fatalf("Unable to create new Flake: %s", err)
	}

	for i := 0; i < b.N; i++ {
		_ = f.NextId()
	}
}
