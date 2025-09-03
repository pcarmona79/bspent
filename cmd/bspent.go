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
// along with bspent. If not, see <https://www.gnu.org/licenses/>.
//

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pcarmona79/bspent/bsp"
)

func toAscii(number uint32) (string, error) {
	bufascii := new(bytes.Buffer)

	err := binary.Write(bufascii, binary.LittleEndian, number)
	if err != nil {
		return "", fmt.Errorf("binary.Write() failed with value %d: %v", number, err)
	}

	return bufascii.String(), nil
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "bspent - BSP entities utility")
	fmt.Fprintln(w, "Usage: bspent <-p|-x> <filename>")
	fmt.Fprintln(w, "  -p: Parse entities inside BSP file.")
	fmt.Fprintln(w, "  -x: Writes to standard output the entities of the BSP file.")
	fmt.Fprintln(w, "  <filename> must be a .bsp or .ent file.")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: Not enough parameters\n")
		usage(os.Stderr)
		os.Exit(-1)
	}

	// operation modes
	extractMode := false
	parseMode := true
	noerrors := true

	if os.Args[1] == "-x" {
		extractMode = true
		parseMode = false
	}

	for _, file := range os.Args[2:] {
		var err error
		var bspfile bsp.BspFile
		bspfile.Filename = file

		// Verify if the file has a .bsp or .ent extension.
		// The .ent files go directly to entity check. The .bsp file
		// first must be readed to get the header and then extract
		// the entities text
		extension := filepath.Ext(file)

		switch extension {
		case ".bsp":
			// read header and extract entities
			err = bspfile.ReadHeader()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error when reading header of '%v': %v\n", file, err)
				noerrors = false
				continue
			}

			if !extractMode {
				fmt.Fprintf(os.Stderr, "Loaded header of file '%v':\n", file)

				magic, _ := toAscii(bspfile.Header.Magic)
				fmt.Fprintf(os.Stderr, "Header magic: %s\n", magic)
				fmt.Fprintf(os.Stderr, "Header version: %d\n", bspfile.Header.Version)

				// show lumps
				for i, lmp := range bspfile.Header.Lump {
					fmt.Fprintf(os.Stderr, "Lump[%d] at offset %d is %d bytes length\n", i,
						lmp.Offset, lmp.Length)
				}

				fmt.Fprintf(os.Stderr, "Reading %d bytes from offset %d to get entities...\n",
					bspfile.Header.Lump[0].Length,
					bspfile.Header.Lump[0].Offset)
			}

			err = bspfile.ReadEntities()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error when reading entities of '%v': %v\n", file, err)
				noerrors = false
				continue
			}

			if extractMode {
				// dump entities to standard output
				bspfile.Entities.WriteEntities(os.Stdout)
			} else if parseMode {
				fmt.Fprintf(os.Stderr, "Validating file '%v'...\n", file)
				bspfile.Entities.Parse(int(bspfile.Header.Lump[0].Length))
				bspfile.Entities.WriteParsed(os.Stdout)
			}

		case ".ent":
			if extractMode {
				fmt.Fprintf(os.Stderr, "Extract mode is not available for '%v'\n", file)
				noerrors = false
				continue
			}

			// load entities file
			var filesize int64
			filesize, err = bspfile.ReadEntitiesFile()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error when reading entities of '%v': %v\n", file, err)
				noerrors = false
				continue
			}

			if parseMode {
				fmt.Fprintf(os.Stderr, "Validating file '%v'...\n", file)
				bspfile.Entities.Parse(int(filesize))
				bspfile.Entities.WriteParsed(os.Stdout)
			}

		default:
			fmt.Fprintf(os.Stderr, "Error: Unsupported format in file '%v'\n", file)
			noerrors = false
			continue
		}
	}

	if !noerrors {
		os.Exit(-1)
	}
	os.Exit(0)
}
