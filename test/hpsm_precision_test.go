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

/**go test --cover hpsm_precision_test.go*/
import (
	"testing"

	proc "scanoss.com/hpsm/pkg"
)

func TestDetectsAccurateStart(t *testing.T) {

	local := "This is line 1\nThis is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n,This is line 5\n"
	remote := "This is line 1\nThis is line 2\nThis is line 3\nThis is line 4\n,This is line 5\nThis is line 6\nThis is line 7\n"
	hashLocal := proc.GetLineHashesFromSource(local)
	hashRemote := proc.GetLineHashesFromSource(remote)
	r := proc.Compare(hashLocal, hashRemote, 5)
	check := (r[0].LStart == 1 && r[0].RStart == 0 && r[0].REnd == 4 && r[0].LEnd == 5)

	if !check {
		t.Errorf("Expected: %v", r)
	}

	local = "This is line 1\n" + local
	hashLocal = proc.GetLineHashesFromSource(local)
	r = proc.Compare(hashLocal, hashRemote, 5)
	check = (r[0].LStart == 2 && r[0].RStart == 0 && r[0].REnd == 4 && r[0].LEnd == 6)

	if !check {
		t.Errorf("Expected: %v", r)
	}

}

func TestDetectsDuplicateStart(t *testing.T) {

	local := "Line0\nLine1\nLine2\nLine3\nLine4\nline5\nLine6\nLine7"
	remote := "Line0\nLine1\nLine2\nLine3\nlinexx\nlineyy\nLine0\nLine1\nLine2\nLine3\nline4"
	hashLocal := proc.GetLineHashesFromSource(local)
	hashRemote := proc.GetLineHashesFromSource(remote)
	r := proc.Compare(hashLocal, hashRemote, 5)
	if len(r) > 0 {
		check := r[0].REnd-r[0].RStart == 3

		if !check {
			t.Errorf("%v", r)
		}
	}

}
func TestDetectsLongestSnippet(t *testing.T) {
	snippet := "This is line 0\nThis is line 1\nThis is line 2\nThis is line 3\nThis is line 4\nThis is line 5\n"
	local := snippet + "This is a middle line\n" + snippet + "This is an additional line\n" + "This is another line\n" + snippet + "\nThis is the end line\n"
	remote := "This is line M\n" + snippet + "This small middle line\n" + snippet + "This is an additional line\n" + "This line sucks\n"
	hashLocal := proc.GetLineHashesFromSource(local)
	hashRemote := proc.GetLineHashesFromSource(remote)
	r := proc.Compare(hashLocal, hashRemote, 3)
	_ = r
	//t.Error(r) //"Expected 2 matches")

}

func TestInverseSnippetOrder(t *testing.T) {
	snippet1 := "sn10\nsn11\nsn12\nsn13\nfinsn1\n"
	snippet2 := "sn20\nsn21\nsn22\nsn23\nfinsn2\n"

	local := "Line1\nLine2\n" + snippet1 + "line3\nline4\n" + snippet2
	//local "Line1\nLine2\nsn10\nsn11\nsn12\nsn13\nfinsn1\nline3\nline4\nsn20\nsn21\nsn22\nsn23\nfinsn2\n
	//remote Line0\sn20\nsn21\nsn22\nsn23\nfinsn2\nsn10\nsn11\nsn12\nsn13\nfinsn1\n
	remote := "Line0\n" + snippet2 + snippet1
	hashLocal := proc.GetLineHashesFromSource(local)
	hashRemote := proc.GetLineHashesFromSource(remote)
	r := proc.Compare(hashLocal, hashRemote, 5)
	got := len(r)
	if got != 2 {
		t.Errorf("Expected: %v, got: %v", 2, r)
	}

}
