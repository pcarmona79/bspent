//
// This file is part of bspent.

// bspent is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at
// your option) any later version.

// bspent is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License
// for more details.

// You should have received a copy of the GNU General Public License
// along with Foobar. If not, see <https://www.gnu.org/licenses/>.
//

package ent

import (
	"fmt"
	"io"
	"os"
)

type Property struct {
	Name string
	Value string
}

type Entity struct {
	Props []Property
}

type Entities struct {
	EntString string
	Ents []Entity
}

func skipComment(text string, line int, idx int) (int, int) {
	var i int

	// comments starts with a \/\/ and ends with an EOL
	for i = idx; i < idx + len(text); i++ {
		c := text[i]
		if c == '\n' {
			line++
			break
		}
	}

	return line, i
}

func skipBlanks(text string, stop byte, idx int) int {
	var i int
	for i = idx; i < idx + 128; i++ {
		c := text[i]
		if c == stop {
			break
		}
	}
	return i
}

func parseKeyValue(text string, line int, start int, end int) (int, int, string, string) {
	var i int
	ikey := 0
	ivalue := 0
	key := ""
	value := ""

	for i = start; i < end; i++ {
		c := text[i]
		// fmt.Printf("parseKeyValue: i=%v c=%c\n", i, c)

		if c == '"' {
			if ikey <= 0 && ivalue <= 0 {
				// found start of key
				i++
				ikey = i
				// fmt.Printf("parseKeyValue: found start of key at %v\n", ikey)

			} else if ikey > 0 && ivalue <= 0 {
				// found end of key
				key = text[ikey:i]
				// fmt.Printf("parseKeyValue: found end of key at %v\n", i)
				// fmt.Printf("parseKeyValue: key is \"%v\"\n", key)
				// skip blanks until next double quote
				i++
				i = skipBlanks(text, '"', i)
				i++
				ivalue = i
				// fmt.Printf("parseKeyValue: skipping blanks until %v\n", i)
				// fmt.Printf("parseKeyValue: found start of value at %v\n", ivalue)

			} else if ikey > 0 && ivalue > 0 {
				// found end of value
				value = text[ivalue:i]
				// fmt.Printf("parseKeyValue: found end of value at %v\n", i)
				// fmt.Printf("parseKeyValue: value is \"%v\"\n", value)
				i++
				break
			}
		}
		if c == '\n' {
			line++
		}
		// fmt.Printf("parseKeyValue: ending loop with i=%d\n", i)
	}

	return line, i, key, value
}

func locateEnt(text string, idx int) int {
	var i int

	// objects starts with a \{ and ends with a \}
	for i = idx; i < idx + len(text); i++ {
		c := text[i]
		if c == '{' {
			continue
		} else if c == '}' {
			break
		}
	}

	return i
}

func (ent *Entities) parseEnt(line int, idx int) (int, int) {
	var key string
	var value string

	entity := Entity{}

	end := locateEnt(ent.EntString, idx)
	for {
		line, idx, key, value = parseKeyValue(ent.EntString, line, idx, end)
		if len(key) <= 0 || len(value) <= 0 {
			break;
		}
		entity.Props = append(entity.Props, Property{key, value})
	}

	ent.Ents = append(ent.Ents, entity)

	return line, idx
}

func (entity *Entity) write(w io.Writer) {
	fmt.Fprintln(w, "{")
	for _, e := range entity.Props {
		fmt.Fprintf(w, "  \"%v\" \"%v\"\n", e.Name, e.Value)
	}
	fmt.Fprintln(w, "}")
}

func (ent *Entities) WriteParsed(w io.Writer) {
	// writes entities in MAP format
	for _, e := range ent.Ents {
		e.write(w)
	}
}

func (ent *Entities) Parse(length int) {
	line := 1

	fmt.Fprintf(os.Stderr, "Parsing %d bytes of entities\n", length)

	// all comments start with "//", whitespaces does not count, objets are surrounded by { and },
	// and inside there are lines key-value, key and value are doble-quoted separately with
	// a whitespace between them
	for i := 0; i < length; i++ {
		c := ent.EntString[i];

		// search for comments
		if c == '/' {
			if i + 1 < length && ent.EntString[i + 1] == '/' {
				fmt.Fprintf(os.Stderr, "Comment at line %v\n", line)
				line, i = skipComment(ent.EntString, line, i)
				continue
			}
		}

		// search for bad comments
		if c == ';' {
			fmt.Fprintf(os.Stderr, "Bad comment at line %v\n", line)
			line, i = skipComment(ent.EntString, line, i)
			continue
		}

		// search for bad comments
		if c == '#' {
			fmt.Fprintf(os.Stderr, "Bad comment at line %v\n", line)
			line, i = skipComment(ent.EntString, line, i)
			continue
		}

		if c == '{' {
			fmt.Fprintf(os.Stderr, "Entity found at line %v\n", line)
			line, i = ent.parseEnt(line, i)
			continue
		}

		if c == '\n' {
			line++
		}
	}
}

func (ent *Entities) Sanitize() string {
	length := len(ent.EntString)

	// sometimes there is a null character at the end
	if ent.EntString[length - 1] == 0 {
		return ent.EntString[0:length - 1]
	}

	return ent.EntString
}

func (ent *Entities) WriteEntities(w io.Writer) (int, error) {
    return w.Write([]byte(ent.Sanitize()))
}
