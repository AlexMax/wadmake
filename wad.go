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

// WadType designates the type of WAD file a WAD is.
type WadType int

const (
	// WadTypeIWAD designates an IWAD, a WAD file that is considered
	// the primary resource file for the game.
	WadTypeIWAD WadType = iota

	// WadTypePWAD designates a PWAD, or a WAD file that is "patched"
	// on top of other resources.
	WadTypePWAD
)

// Wad is a stucture that contains all the data necessary to marshall
// a WAD file.
type Wad struct {
	WadType WadType
	Lumps   Directory
}

// NewWad creates a new Wad structure.
func NewWad(wadType WadType) *Wad {
	wad := &Wad{}
	wad.WadType = wadType
	wad.Lumps = Directory{}

	return wad
}

// Decode decodes WAD file data passed into the reader into a Wad
// structure.
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

// Encode encodes the passed Wad structure into WAD file data.
func Encode(w io.Writer, wad *Wad) error {
	const maxInt32 = 2147483647
	const headerOffset = 12

	var n int
	var err error

	// Write header
	switch wad.WadType {
	case WadTypeIWAD:
		n, err = w.Write([]byte("IWAD"))
	case WadTypePWAD:
		n, err = w.Write([]byte("PWAD"))
	default:
		return errors.New("could not write header of unknown wad type")
	}

	if err != nil {
		return err
	} else if n < 4 {
		return errors.New("could not write header")
	}

	var alldata bytes.Buffer
	var infotable bytes.Buffer

	for i := 0; i < len(wad.Lumps); i++ {
		// Write lump position
		alldatapos := int64(alldata.Len()) + headerOffset
		if alldatapos > maxInt32 {
			return errors.New("could not write lump position")
		}
		err = binary.Write(&infotable, binary.LittleEndian, int32(alldatapos))
		if err != nil {
			return err
		}

		// Write lump data
		n, err = alldata.Write(wad.Lumps[i].Data)
		if err != nil {
			return err
		} else if n < len(wad.Lumps[i].Data) {
			return errors.New("could not write lump data")
		}

		// Write lump size
		dataSize := int64(len(wad.Lumps[i].Data))
		if dataSize > maxInt32 {
			return errors.New("could not write lump size")
		}
		err := binary.Write(&infotable, binary.LittleEndian, int32(dataSize))
		if err != nil {
			return err
		}

		// Write lump name.  Lump names are a maximum of 8 characters,
		// and any shorter names must end with a null terminator.
		if len(wad.Lumps[i].Name) > 8 {
			return errors.New("lump name is too long")
		}
		var namebuffer [8]byte
		n := copy(namebuffer[:], wad.Lumps[i].Name)
		if n < len(wad.Lumps[i].Name) {
			return errors.New("could not copy lump name")
		}
		n, err = infotable.Write(namebuffer[:])
		if err != nil {
			return err
		} else if n < len(namebuffer) {
			return errors.New("could not write lump name")
		}
	}

	// Write number of lumps
	if len(wad.Lumps) > maxInt32 {
		return errors.New("too many lumps")
	}
	err = binary.Write(w, binary.LittleEndian, int32(len(wad.Lumps)))
	if err != nil {
		return err
	}

	// Write offset of infotable
	alldatapos := int64(alldata.Len()) + headerOffset
	if alldatapos > maxInt32 {
		return errors.New("invalid infotable offset")
	}
	err = binary.Write(w, binary.LittleEndian, int32(alldatapos))
	if err != nil {
		return err
	}

	// Write data
	n, err = w.Write(alldata.Bytes())
	if err != nil {
		return err
	} else if n < alldata.Len() {
		return errors.New("could not write data")
	}

	// Write infotable
	n, err = w.Write(infotable.Bytes())
	if err != nil {
		return err
	} else if n < infotable.Len() {
		return errors.New("could not write infotable")
	}

	return nil
}
