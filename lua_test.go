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

func TestNewLuaEnvironment(t *testing.T) {
	// Environment should be created correctly
	l := NewLuaEnvironment()

	// Don't leave anything on the stack when creating the environment
	if l.Top() != 0 {
		t.Fatalf("environment does not have a clean stack (%d)", l.Top())
	}

	// wad package should exist
	l.Global("wad")
	if l.TypeOf(-1) != lua.TypeTable {
		t.Fatal("wad package does not exist")
	}
}
