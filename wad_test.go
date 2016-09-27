/*
 *  Copyright 2016 Alex Mayfield
 *
 *  This file is part of WADmake.
 *
 *  WADmake is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  WADmake is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with WADmake.  If not, see <http://www.gnu.org/licenses/>.
 */

package wadmake

import (
	"os"
	"testing"
)

func TestWadDecode(t *testing.T) {
	file, err := os.Open("wadmake_test.wad")
	if err != nil {
		t.Fatal(err.Error())
	}

	wad, err := Decode(file)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Contains 11 lumps
	if len(wad.Lumps) != 11 {
		t.Error("incorrect lump count in decoded WAD file")
	}

	// MAP01 is < 8 characters
	if wad.Lumps[0].Name != "MAP01" {
		t.Error("incorrect short lump name in decoded WAD file")
	}

	// BLOCKMAP is 8 characters
	if wad.Lumps[10].Name != "BLOCKMAP" {
		t.Error("incorrect lump name in decoded WAD file")
	}

	// THINGS lump is 10 bytes per THING
	if len(wad.Lumps[1].Data)%10 != 0 {
		t.Error("incorrect lump data length in decoded WAD file")
	}
}
