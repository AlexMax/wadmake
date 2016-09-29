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

// Create empty Lumps userdata
func wadCreateLumps(l *lua.State) int {
	l.PushUserData(&Directory{})
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
	l.PushUserData(&wad.Lumps)
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
	{"find", lumpsFind},
	{"get", lumpsGet},
	{"insert", lumpsInsert},
	{"packwad", lumpsPackWAD},
	{"packzip", nil},
	{"remove", lumpsRemove},
	{"set", lumpsSet},
	{"writewad", lumpsWriteWAD},
	{"writezip", nil},
	{"__len", lumpsLen},
	{"__tostring", lumpsToString},
}

// Checks for Lumps userdata at a specific stack index, usually 1.
func checkLumps(l *lua.State, index int) *Directory {
	data, ok := lua.CheckUserData(l, index, lumpsHandle).(*Directory)
	if !ok {
		lua.Errorf(l, "type asserion failed")
	} else if data == nil {
		lua.Errorf(l, "nil pointer")
	}

	return data
}

// Finds a lump by name, optionally starting in the middle.  Returns
// the location of the lump, or nil if not found.
func lumpsFind(l *lua.State) int {
	data := checkLumps(l, 1)
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
		start = len(*data) + start + 1
	}

	if start > len(*data) {
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

// Returns a lump name and raw data, or nil if nothing was found.
func lumpsGet(l *lua.State) int {
	data := checkLumps(l, 1)
	index := lua.CheckInteger(l, 2)
	if index < 1 || index > len(*data) {
		l.PushNil()
		return 1
	}

	l.PushString((*data)[index-1].Name)
	l.PushString(string((*data)[index-1].Data))

	return 2
}

// Insert lump data into the directory.
func lumpsInsert(l *lua.State) int {
	data := checkLumps(l, 1)

	if l.IsNumber(2) {
		// Second parameter is index to push to
		index := lua.CheckInteger(l, 2)
		if index < 1 || index > len(*data) {
			lua.ArgumentError(l, 2, "index out of range")
		}

		lump := Lump{
			Name: lua.CheckString(l, 3),
			Data: []byte(lua.CheckString(l, 4)),
		}

		*data = append(*data, Lump{})
		copy((*data)[index:], (*data)[index-1:])
		(*data)[index-1] = lump
	} else {
		// Append to end
		lump := Lump{
			Name: lua.CheckString(l, 2),
			Data: []byte(lua.CheckString(l, 3)),
		}
		*data = append(*data, lump)
	}

	return 0
}

// Pack WAD file into string.
func lumpsPackWAD(l *lua.State) int {
	data := checkLumps(l, 1)

	// Create WAD structure
	wad := NewWad(WadTypePWAD)
	wad.Lumps = *data

	// Write into bytebuffer
	buffer := bytes.Buffer{}
	err := Encode(&buffer, wad)
	if err != nil {
		lua.Errorf(l, "could not encode data (%s)", err.Error())
	}

	// Return bytebuffer
	l.PushString(buffer.String())

	return 1
}

// Remove lump data from the directory.
func lumpsRemove(l *lua.State) int {
	data := checkLumps(l, 1)

	// Second parameter is index to remove from
	index := lua.CheckInteger(l, 2)
	if index < 1 || index > len(*data) {
		lua.ArgumentError(l, 2, "index out of range")
	}

	*data = append((*data)[:index-1], (*data)[index:]...)

	return 0
}

// Set lump data at a specific point in the directory.  You can omit
// the name or data but not both.
func lumpsSet(l *lua.State) int {
	data := checkLumps(l, 1)

	// Second parameter is index to set
	index := lua.CheckInteger(l, 2)
	if index < 1 || index > len(*data) {
		lua.ArgumentError(l, 2, "index out of range")
	}

	nametype := l.TypeOf(3)
	datatype := l.TypeOf(4)

	// Name is string if present, nil if not
	if nametype != lua.TypeString && nametype != lua.TypeNil {
		lua.ArgumentError(l, 3, "must be string or nil")
	}

	// Data is string if present, nil if not
	if datatype != lua.TypeString && datatype != lua.TypeNone {
		lua.ArgumentError(l, 4, "must be string, if present")
	}

	if nametype == lua.TypeNil || datatype == lua.TypeNone {
		// If one of the parameters is missing, we need the original
		if nametype == lua.TypeString {
			lname, _ := l.ToString(3)
			(*data)[index-1].Name = lname
		}

		if datatype == lua.TypeString {
			ldata, _ := l.ToString(4)
			(*data)[index-1].Data = []byte(ldata)
		}
	} else {
		// Both parameters, so a brand new lump.
		lname, _ := l.ToString(3)
		ldata, _ := l.ToString(4)

		lump := Lump{
			Name: lname,
			Data: []byte(ldata),
		}

		(*data)[index-1] = lump
	}

	return 0
}

// Write WAD file to disk.
func lumpsWriteWAD(l *lua.State) int {
	data := checkLumps(l, 1)
	filename := lua.CheckString(l, 2)

	// Create WAD structure
	wad := NewWad(WadTypePWAD)
	wad.Lumps = *data

	// Open file
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		lua.Errorf(l, "could not open file (%s)", err.Error())
	}

	// Write into file
	err = Encode(file, wad)
	if err != nil {
		lua.Errorf(l, "could not encode data (%s)", err.Error())
	}

	return 0
}

// Length of directory.
func lumpsLen(l *lua.State) int {
	data := checkLumps(l, 1)
	l.PushInteger(len(*data))
	return 1
}

// A nice way of printing the Lumps userdata.
func lumpsToString(l *lua.State) int {
	data := checkLumps(l, 1)

	dataLen := len(*data)
	if dataLen == 1 {
		l.PushFString("%s: %p, %d lump", lumpsHandle, data, 1)
	} else {
		l.PushFString("%s: %p, %d lumps", lumpsHandle, data, dataLen)
	}
	return 1
}

// WadLumpsOpen adds all lump-related functions to the table located at
// the top of the stack of the pased lua state.
func WadLumpsOpen(l *lua.State) error {
	// [wadlib] Set global wad functions.
	lua.SetFunctions(l, wadMethods, 0)

	// Create the Lumps userdata with associated functions.
	ok := lua.NewMetaTable(l, lumpsHandle)
	if !ok {
		return errors.New("could not create Lumps metatable")
	}

	// [wadlib][Lumpsmeta]
	l.PushValue(-1)
	// [wadlib][Lumpsmeta][Lumpsmeta]
	l.SetField(-2, "__index")
	// [wadlib][Lumpsmeta]
	lua.SetFunctions(l, lumpsMethods, 0)
	l.Pop(1)
	// [wadlib]

	return nil
}
