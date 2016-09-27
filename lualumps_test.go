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
	"testing"

	lua "github.com/Shopify/go-lua"
)

// Lumps can be created from scratch
func TestCreateLumps(t *testing.T) {
	l := NewLuaEnvironment()

	// Test the function invocation
	err := lua.DoString(l, "return wad.createLumps()")
	if err != nil {
		t.Fatal(err.Error())
	}

	// A fresh Lumps must be of userdata type Lumps
	_, ok := lua.CheckUserData(l, -1, "Lumps").(*Directory)
	if ok == false {
		t.Fatal("Lumps is not *Directory")
	}

	// A fresh Lumps must be empty
	if lua.LengthEx(l, -1) != 0 {
		t.Fatal("Lumps is not empty")
	}
}

// Lumps can be read from a file
func TestReadWAD(t *testing.T) {
	l := NewLuaEnvironment()

	// Test the function invocation
	err := lua.DoString(l, "return wad.readwad('wadmake_test.wad')")
	if err != nil {
		t.Fatal(err.Error())
	}

	// Make sure we've actually got two items in valid locations on
	// the stack.
	if l.Top() != 2 {
		t.Fatal("incorrect stack size")
	}

	// First parameter must be correct WAD type
	wadType := lua.CheckString(l, -1)
	if wadType != "pwad" {
		t.Fatal("wad type is not pwad")
	}

	// Second parameter must be of userdata type Lumps
	_, ok := lua.CheckUserData(l, -2, "Lumps").(*Directory)
	if ok == false {
		t.Fatal("Lumps is not *Directory")
	}

	// Second parameter must contain the correct number of lumps
	if lua.LengthEx(l, -2) != 11 {
		t.Fatal("Lumps is not empty")
	}
}

func readWad(t *testing.T) *lua.State {
	l := NewLuaEnvironment()

	err := lua.DoString(l, "lumps = wad.readwad('wadmake_test.wad')")
	if err != nil {
		t.Fatal(err.Error())
	}

	return l
}

// Find a lump by index
func TestLumpsFind(t *testing.T) {
	l := readWad(t)
	lua.DoString(l, "return lumps:find('SIDEDEFS')")

	if l.Top() != 1 {
		t.Fatal("incorrect stack size")
	}

	if lua.CheckInteger(l, -1) != 4 {
		t.Error("incorrect lump position")
	}
}

// Find a lump index given a starting index
func TestLumpsFindIndex(t *testing.T) {
	l := readWad(t)
	lua.DoString(l, "return lumps:find('SIDEDEFS', 2)")

	if l.Top() != 1 {
		t.Fatalf("incorrect stack size")
	}

	if lua.CheckInteger(l, -1) != 4 {
		t.Error("incorrect lump position")
	}
}

// Find a lump index given a negative starting index
func TestLumpsFindNegativeIndex(t *testing.T) {
	l := readWad(t)
	lua.DoString(l, "return lumps:find('SIDEDEFS', -9)")

	if l.Top() != 1 {
		t.Fatal("incorrect stack size")
	}

	if lua.CheckInteger(l, -1) != 4 {
		t.Error("incorrect lump position")
	}
}

// Return a lump name and data
func TestLumpsGet(t *testing.T) {
	l := readWad(t)
	lua.DoString(l, "return lumps:get(2)")

	if l.Top() != 2 {
		t.Fatal("incorrect stack size")
	}

	if lua.CheckString(l, -2) != "THINGS" {
		t.Error("incorrect lump name")
	}

	if len(lua.CheckString(l, -1))%10 != 0 {
		t.Error("incorrect lump data length")
	}
}

// Insert a lump into the middle
func TestLumpsInsert(t *testing.T) {
	l := readWad(t)
	lua.DoString(l, "lumps:insert(2, 'TEST', 'hissy');return lumps:get(2)")

	if l.Top() != 2 {
		t.Fatal("incorrect stack size")
	}

	if lua.CheckString(l, -2) != "TEST" {
		t.Error("incorrect lump name")
	}

	if lua.CheckString(l, -1) != "hissy" {
		t.Error("incorrect lump data")
	}

	lua.DoString(l, "return lumps")

	if lua.LengthEx(l, -1) != 12 {
		t.Error("incorrect lump count")
	}
}

// Append a lump to the end
func TestLumpsInsertAppend(t *testing.T) {
	l := readWad(t)
	lua.DoString(l, "lumps:insert('MAP02', 'hissy');return lumps:get(12)")

	if l.Top() != 2 {
		t.Fatal("incorrect stack size")
	}

	if lua.CheckString(l, -2) != "MAP02" {
		t.Error("incorrect lump name")
	}

	if lua.CheckString(l, -1) != "hissy" {
		t.Error("incorrect lump data")
	}

	lua.DoString(l, "return lumps")

	if lua.LengthEx(l, -1) != 12 {
		t.Error("incorrect lump count")
	}
}
