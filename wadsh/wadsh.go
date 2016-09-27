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

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/AlexMax/wadmake"
	"github.com/Shopify/go-lua"
	"github.com/chzyer/readline"
)

func main() {
	env := wadmake.NewLuaEnvironment()

	if len(os.Args) > 1 {
		if os.Args[1] == "-" {
			// First parameter is -, read script from stdin.
		} else {
			// Read scripts in as filenames.
			err := lua.DoFile(env, os.Args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	rl, err := readline.New("> ")
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Fprint(os.Stderr, "WADmake shell\n")
	fmt.Fprint(os.Stderr, "Press Ctrl-C to quit the shell.\n")
	fmt.Fprint(os.Stderr, "Press Ctrl-D on an empty line to quit the shell.\n")

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}

		err = lua.DoString(env, line)
		if err != nil {
			errString, ok := env.ToString(-1)
			if ok {
				fmt.Fprintf(os.Stderr, "%s\n", errString)
			} else {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		}
	}
}
