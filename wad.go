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
	"encoding/binary"
	"errors"
	"io"
)

type WadType int

const (
	WadTypeIWAD WadType = iota
	WadTypePWAD
)

type Wad struct {
	WadType WadType
	Lumps   Directory
}

func NewWad(wadType WadType) *Wad {
	wad := &Wad{}
	wad.WadType = wadType
	wad.Lumps = Directory{}

	return wad
}

func Decode(r io.ReadSeeker) (*Wad, error) {
	// WAD identifier
	var identifier [4]byte
	n, err := r.Read(identifier[:])
	if err != nil {
		return nil, err
	} else if n < 4 {
		return nil, errors.New("could not read WAD identifier")
	}

	var wad *Wad
	switch identifier {
	case [...]byte{'I', 'W', 'A', 'D'}:
		wad = NewWad(WadTypeIWAD)
	case [...]byte{'P', 'W', 'A', 'D'}:
		wad = NewWad(WadTypePWAD)
	default:
		return nil, errors.New("invalid WAD identifier")
	}

	// Number of lumps
	var numlumps int32
	err = binary.Read(r, binary.LittleEndian, &numlumps)
	if err != nil {
		return nil, err
	} else if numlumps < 0 {
		return nil, errors.New("too many lumps")
	}

	// Infotable location
	var infotablefs int32
	err = binary.Read(r, binary.LittleEndian, &infotablefs)
	if infotablefs < 0 {
		return nil, errors.New("infotable out of range")
	}
	_, err = r.Seek(int64(infotablefs), io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read infotable
	var i int32
	for i = 0; i < numlumps; i++ {
		// Read file position and size
		var filepos, size int32
		err = binary.Read(r, binary.LittleEndian, &filepos)
		if err != nil {
			return nil, err
		}
		err = binary.Read(r, binary.LittleEndian, &size)
		if err != nil {
			return nil, err
		}

		// Read name
		var name [8]byte
		n, err := r.Read(name[:])
		if err != nil {
			return nil, err
		} else if n < 8 {
			return nil, errors.New("could not read name")
		}

		// Create lump
		lump := Lump{}

		// Name is either null-terminated and less than 8 bytes, or
		// 8 bytes exactly and not null-terminated.
		index := bytes.IndexByte(name[:], byte(0))
		if index != -1 {
			lump.Name = string(name[:index])
		} else {
			lump.Name = string(name[:])
		}

		// If size is 0, file position could be nonsense, so only
		// attempt to read data if size > 0.
		if size > 0 {
			if filepos < 0 {
				return nil, errors.New("filepos out of range")
			}

			// Store current position
			info, err := r.Seek(0, io.SeekCurrent)
			if err != nil {
				return nil, err
			}

			// Read data.
			_, err = r.Seek(int64(filepos), io.SeekStart)
			if err != nil {
				return nil, err
			}

			var data = make([]byte, size)
			n, err = r.Read(data)
			if err != nil {
				return nil, err
			} else if int32(n) < size {
				return nil, errors.New("could not read data")
			}
			lump.Data = data

			// Go back to old position
			_, err = r.Seek(info, io.SeekStart)
			if err != nil {
				return nil, err
			}
		}

		wad.Lumps = append(wad.Lumps, lump)
	}

	return wad, nil
}

func Encode(w io.Writer, wad *Wad) error {
	return nil
}
