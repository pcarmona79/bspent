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

package bsp

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/pcarmona79/bspent/ent"
)

type BspLump struct {
	Offset uint32
	Length uint32
}

type BspHeader struct {
	Magic uint32
	Version uint32
	Lump [19]BspLump
}

type BspFile struct {
	Filename string
	Header BspHeader
	Entities ent.Entities
}

func (bsp *BspFile) ReadHeader() error {
	file, err := os.Open(bsp.Filename)
	if err != nil {
		return fmt.Errorf("Error opening file '%v' to read header\n", bsp.Filename)
	}
	defer file.Close()

	err = binary.Read(file, binary.LittleEndian, &bsp.Header)
	if err != nil {
		return fmt.Errorf("Error reading header: %v\n", err)
	}

	return nil
}

func (bsp *BspFile) loadEntities(file *os.File, bytesNeeded int64) error {
	var bytesRead int
	var err error
	var finalbuf []byte

	rest := bytesNeeded % 1024
	if rest > 0 {
		finalbuf = make([]byte, rest)
	}

	buf := make([]byte, 1024)
	for {
		if bytesNeeded >= 1024 {
			bytesRead, err = file.Read(buf)
		} else {
			bytesRead, err = file.Read(finalbuf)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Error reading entities chunk: %v\n", err)
		}

		if bytesNeeded >= 1024 {
			bsp.Entities.EntString += string(buf)
		} else {
			bsp.Entities.EntString += string(finalbuf)
		}

		bytesNeeded -= int64(bytesRead)
		if bytesNeeded <= 0 {
			break;
		}
	}

	return nil
}

func (bsp *BspFile) ReadEntities() error {
	file, err := os.Open(bsp.Filename)
	if err != nil {
		return fmt.Errorf("Error opening file '%v' to read entities\n", bsp.Filename)
	}
	defer file.Close()

	file.Seek(int64(bsp.Header.Lump[0].Offset), 0)

	return bsp.loadEntities(file, int64(bsp.Header.Lump[0].Length))
}

func (bsp *BspFile) ReadEntitiesFile() (int64, error) {
	file, err := os.Open(bsp.Filename)
	if err != nil {
		return -1, fmt.Errorf("Error opening entities file '%v'\n", bsp.Filename)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return -1, fmt.Errorf("Error getting file information from '%v'\n", bsp.Filename)
	}
	bytesNeeded := info.Size()

	return bytesNeeded, bsp.loadEntities(file, bytesNeeded)
}
