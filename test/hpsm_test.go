package test

import (
	"testing"

	proc "scanoss.com/hpsm/pkg"
)

func TestDetectsStarts(t *testing.T) {
	expected := true
	local := "This is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n,This is line 5\n"
	remote := "This is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n,This is line 5\nThis is line 6\nThis is line 7\n"
	hashLocal := proc.GetLineHashesFromSource(local)
	hashRemote := proc.GetLineHashesFromSource(remote)
	r := proc.Compare(hashLocal, hashRemote, 5)

	got := len(r)
	if got != 1 {
		t.Errorf("Expected: %v, got: %v", expected, got)
	}
}

func TestDetectsEnds(t *testing.T) {
	expected := true
	local := "This is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n,This is line 5\n"
	remote := "This is line -1\nThis is line 0\nThis is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n,This is line 5\nThis is line 6\nThis is line 7\n"
	hashLocal := proc.GetLineHashesFromSource(local)
	hashRemote := proc.GetLineHashesFromSource(remote)
	r := proc.Compare(hashLocal, hashRemote, 5)

	got := len(r)
	if got != 1 {
		t.Errorf("Expected: %v, got: %v", expected, got)
	}

}
func TestDetectsMiddle(t *testing.T) {
	expected := true
	local := "This is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n,This is line 5\n"
	remote := "This is line -1\nThis is line 0\nThis is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n,This is line 5\nThis is line 6\nThis is line 7\nThis is line 8\nThis is line 9"
	hashLocal := proc.GetLineHashesFromSource(local)
	hashRemote := proc.GetLineHashesFromSource(remote)
	r := proc.Compare(hashLocal, hashRemote, 5)

	got := len(r)
	if got != 1 {
		t.Errorf("Expected: %v, got: %v", expected, got)
	}

}
func TestDetectsThreshold(t *testing.T) {
	expected := true
	local := "This is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n,This is line 5\n"
	remote := "This is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n"
	hashLocal := proc.GetLineHashesFromSource(local)
	hashRemote := proc.GetLineHashesFromSource(remote)
	r := proc.Compare(hashLocal, hashRemote, 5)

	got := len(r)
	if got != 0 {
		t.Errorf("Expected: %v, got: %v", expected, got)
	}
	remote += "This is line 5\n"

	hashRemote = proc.GetLineHashesFromSource(remote)
	r = proc.Compare(hashLocal, hashRemote, 5)

	got = len(r)
	if got != 1 {
		t.Errorf("Expected: %v, got: %v", expected, got)
	}

}
