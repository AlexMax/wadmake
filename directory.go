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

type Lump struct {
	Name string
	Data []byte
}

type Directory []Lump

// Search searches for a specific lump by name and returns its position
// and true if found, or zero and false if not found.
func (dir *Directory) Search(name string, start int) (int, bool) {
	for index, lump := range *dir {
		if name == lump.Name {
			return index, true
		}
	}

	return 0, false
}
