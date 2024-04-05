// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2022 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
func TestWithLongLineTrimmed(t *testing.T) {
	expected := true
	local := "This is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n,This is line 5\n"
	remote := "This is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n"
	newVeryLongLine := []byte{}
	for k := 0; k < 4002; k++ {
		newVeryLongLine = append(newVeryLongLine, []byte("hello")...)
	}
	local = local + "\n" + string(newVeryLongLine)

	hashLocal := proc.GetLineHashesFromSource(local)
	hashRemote := proc.GetLineHashesFromSource(remote)
	r := proc.Compare(hashLocal, hashRemote, 4)

	got := len(r)
	if got != 1 {
		t.Errorf("Expected: %v, got: %v", expected, got)
	}

}
