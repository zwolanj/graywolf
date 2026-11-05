//go:build !android && (linux || cgo)

package kiss

import (
	"sort"
	"testing"
)

func TestBLEScanDedup(t *testing.T) {
	d := newBLEScanDedup()

	if !d.shouldReport("a", "") {
		t.Fatal("first sighting must be reported")
	}
	if d.shouldReport("a", "") {
		t.Error("repeat nameless sighting must be dropped")
	}
	if got := d.unnamed(); len(got) != 1 || got[0] != "a" {
		t.Errorf("unnamed = %v, want [a]", got)
	}
	if !d.shouldReport("a", "B.B. Link Atom") {
		t.Error("name arriving late must be reported")
	}
	if d.shouldReport("a", "B.B. Link Atom") {
		t.Error("repeat named sighting must be dropped")
	}
	if d.shouldReport("a", "") {
		t.Error("nameless sighting after a named one must be dropped")
	}
	if got := d.unnamed(); len(got) != 0 {
		t.Errorf("unnamed = %v, want none", got)
	}

	if !d.shouldReport("b", "Named") || !d.shouldReport("c", "") {
		t.Fatal("new devices must be reported")
	}
	got := d.unnamed()
	sort.Strings(got)
	if len(got) != 1 || got[0] != "c" {
		t.Errorf("unnamed = %v, want [c]", got)
	}
}
