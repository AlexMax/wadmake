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
	"fmt"

	"github.com/Shopify/go-lua"
)

func wadOpen(l *lua.State) int {
	l.NewTable()
	WadLumpsOpen(l)

	return 1
}

func NewLuaEnvironment() *lua.State {
	l := lua.NewState()

	lua.OpenLibraries(l)

	lua.Require(l, "wad", wadOpen, true)
	l.Pop(1)

	return l
}

func LuaDebugStack(state *lua.State) {
	top := state.Top()
	fmt.Printf("Top: %d\n", top)
	for i := 1; i <= top; i++ {
		typename := lua.TypeNameOf(state, i)
		meta, _ := lua.ToStringMeta(state, i)
		state.Pop(1)

		fmt.Printf("%s: %s\n", typename, meta)
	}
}
