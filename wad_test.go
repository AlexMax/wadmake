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
	"bytes"
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

func TestWadEncode(t *testing.T) {
	wad := NewWad(WadTypePWAD)

	wad.Lumps = append(wad.Lumps, Lump{
		Name: "TEST",
		Data: []byte("hissy"),
	})
	wad.Lumps = append(wad.Lumps, Lump{
		Name: "TESTTWO",
		Data: []byte("god only knows"),
	})

	var buffer bytes.Buffer
	err := Encode(&buffer, wad)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := []byte{
		// "PWAD"
		0x50, 0x57, 0x41, 0x44,
		// Number of lumps
		0x2, 0x0, 0x0, 0x0,
		// Location of infotable
		0x1f, 0x0, 0x0, 0x0,
		// "hissy" (Data start)
		0x68, 0x69, 0x73, 0x73, 0x79,
		// "god only knows"
		0x67, 0x6f, 0x64, 0x20, 0x6f, 0x6e, 0x6c, 0x79, 0x20, 0x6b, 0x6e, 0x6f, 0x77, 0x73,
		// Lump location (Infotable start)
		0xc, 0x0, 0x0, 0x0,
		// Lump size
		0x5, 0x0, 0x0, 0x0,
		// "TEST"
		0x54, 0x45, 0x53, 0x54, 0x0, 0x0, 0x0, 0x0,
		// Lump location
		0x11, 0x0, 0x0, 0x0,
		// Lump size
		0xe, 0x0, 0x0, 0x0,
		// "TESTTWO"
		0x54, 0x45, 0x53, 0x54, 0x54, 0x57, 0x4f, 0x0,
	}

	if !bytes.Equal(expected, buffer.Bytes()) {
		t.Error("encoded WAD does not match expected")
	}
}
