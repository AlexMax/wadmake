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
	"errors"
	"io"
	"os"
	"strings"

	lua "github.com/Shopify/go-lua"
)

const lumpsHandle = "Lumps"

var wadMethods = []lua.RegistryFunction{
	{"createLumps", wadCreateLumps},
	{"readwad", wadReadWAD},
	{"unpackwad", wadUnpackWAD},
	{"unpackzip", nil},
}

func wadCreateLumps(l *lua.State) int {
	l.PushUserData(Directory{})
	lua.SetMetaTableNamed(l, lumpsHandle)

	return 1
}

// Load WAD file from disk and return the WAD type and lumps
func wadReadWAD(l *lua.State) int {
	// Read WAD data from filename parameter
	filename := lua.CheckString(l, 1)

	file, err := os.Open(filename)
	if err != nil {
		lua.Errorf(l, err.Error())
	}

	return commonUnpackWAD(l, file)
}

// Read WAD file data and return the WAD type and lumps
func wadUnpackWAD(l *lua.State) int {
	// Read WAD data from string parameter
	buffer := lua.CheckString(l, 1)

	return commonUnpackWAD(l, strings.NewReader(buffer))
}

// Common functionality used to unpack WAD from file or buffer
func commonUnpackWAD(l *lua.State, r io.ReadSeeker) int {
	// Decoade WAD data into data
	wad, err := Decode(r)
	if err != nil {
		lua.Errorf(l, err.Error())
	}

	// Lump data
	l.PushUserData(wad.Lumps)
	lua.SetMetaTableNamed(l, lumpsHandle)

	// WAD type
	switch wad.WadType {
	case WadTypeIWAD:
		l.PushString("iwad")
	case WadTypePWAD:
		l.PushString("pwad")
	default:
		lua.Errorf(l, "unknown wad type")
	}

	return 2
}

var lumpsMethods = []lua.RegistryFunction{
	{"find", lumpFind},
	{"get", lumpGet},
	{"insert", nil},
	{"remove", nil},
	{"set", nil},
	{"packwad", nil},
	{"packzip", nil},
	{"__gc", nil},
	{"__len", lumpsLen},
	{"__tostring", lumpsToString},
}

func lumpFind(l *lua.State) int {
	data, _ := lua.CheckUserData(l, 1, lumpsHandle).(Directory)

	name := lua.CheckString(l, 2)

	// Use a start parameter if we pass it, otherwise the default index
	// to use is 1.
	var start int
	if l.TypeOf(3) != lua.TypeNone {
		start = lua.CheckInteger(l, 3)
	} else {
		start = 1
	}

	// Normalize start parameter.
	if start == 0 {
		// Pretend it's one.
		start = 1
	} else if start < 0 {
		// Negative values are from the end.
		start = len(data) + start + 1
	}

	if start > len(data) {
		l.PushNil()
		return 1
	}

	index, ok := data.Search(name, start-1)
	if ok {
		l.PushInteger(index + 1)
	} else {
		l.PushNil()
	}

	return 1
}

func lumpGet(l *lua.State) int {
	data, _ := lua.CheckUserData(l, 1, lumpsHandle).(Directory)

	index := lua.CheckInteger(l, 2)
	if index < 1 || index > len(data) {
		l.PushNil()
		return 1
	}

	l.PushString(data[index-1].Name)
	l.PushString(string(data[index-1].Data))

	return 2
}

func lumpsLen(l *lua.State) int {
	data, _ := lua.CheckUserData(l, 1, lumpsHandle).(Directory)
	l.PushInteger(len(data))

	return 1
}

func lumpsToString(l *lua.State) int {
	data, _ := lua.CheckUserData(l, 1, lumpsHandle).(Directory)
	if len(data) == 1 {
		l.PushFString("%s: %p, %d lump", lumpsHandle, data, 1)
	} else {
		l.PushFString("%s: %p, %d lumps", lumpsHandle, data, len(data))
	}
	return 1
}

func WadLumpsOpen(state *lua.State) error {
	// [wadlib] Set global wad functions.
	lua.SetFunctions(state, wadMethods, 0)

	// Create the Lumps userdata with associated functions.
	ok := lua.NewMetaTable(state, lumpsHandle)
	if !ok {
		return errors.New("could not create Lumps metatable")
	}

	// [wadlib][Lumpsmeta]
	state.PushValue(-1)
	// [wadlib][Lumpsmeta][Lumpsmeta]
	state.SetField(-2, "__index")
	// [wadlib][Lumpsmeta]
	lua.SetFunctions(state, lumpsMethods, 0)
	state.Pop(1)
	// [wadlib]

	return nil
}
