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
	_, ok := lua.CheckUserData(l, -1, "Lumps").(Directory)
	if ok == false {
		t.Fatal("Lumps is not Directory")
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

	// First parameter must be correct WAD type
	wadType := lua.CheckString(l, -1)
	if wadType != "pwad" {
		t.Fatal("wad type is not pwad")
	}

	// Second parameter must be of userdata type Lumps
	_, ok := lua.CheckUserData(l, -2, "Lumps").(Directory)
	if ok == false {
		t.Fatal("Lumps is not Directory")
	}

	// Second parameter must contain the correct number of lumps
	if lua.LengthEx(l, -2) != 11 {
		t.Fatal("Lumps is not empty")
	}
}
