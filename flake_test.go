package flake

import (
	"sort"
	"testing"
)

func TestNewFlake(t *testing.T) {
	f, err := New()

	if err != nil {
		t.Errorf("Unable to create new ID generator: %s", err)
	}

	var ids []string

	for i := 0; i < 4; i++ {
		id := f.NextId()

		ids = append(ids, id.String())
	}

	if !sort.StringsAreSorted(ids) {
		t.Errorf("IDs are not k-sorted!")
	}
}
