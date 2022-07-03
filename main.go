package main

import (
	"os"
	"strconv"

	. "github.com/stevegt/goadapt"
	"github.com/stevegt/godevmem/devmem"
)

/*
 * godevmem: Go CLI and library to read/write from/to any location in memory.
 *
 *  Copyright (C) 2022, Steve Traugott <stevegt@t7a.org>
 *
 * - Translated to Go from the C version at https://github.com/rcn-ee/devmem2
 * - The following is the original copyright notice from the C version:
 *
 **************************************************************************
 * devmem2.c: Simple program to read/write from/to any location in memory.
 *
 *  Copyright (C) 2000, Jan-Derk Bakker (J.D.Bakker@its.tudelft.nl)
 *
 *
 * This software has been developed for the LART computing board
 * (http://www.lart.tudelft.nl/). The development has been sponsored by
 * the Mobile MultiMedia Communications (http://www.mmc.tudelft.nl/)
 * and Ubiquitous Communications (http://www.ubicom.tudelft.nl/)
 * projects.
 *
 * The author can be reached at:
 *
 *  Jan-Derk Bakker
 *  Information and Communication Theory Group
 *  Faculty of Information Technology and Systems
 *  Delft University of Technology
 *  P.O. Box 5031
 *  2600 GA Delft
 *  The Netherlands
 *
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 *
 */

var usage = `
Usage: %s { address } [ size [ data ] ]
	address : memory address to act upon
	size    : data bit width (default 32): 8, 16, 32, 64
	data    : data to be written
`

func main() {

	if len(os.Args) < 2 {
		Fpf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}

	target, err := strconv.ParseInt(os.Args[1], 0, 64)
	Ck(err)

	var size int64 = 32
	if len(os.Args) > 2 {
		size, err = strconv.ParseInt(os.Args[2], 0, 64)
		Ck(err)
	}

	var newVal uint64
	write := false
	if len(os.Args) > 3 {
		newVal, err = strconv.ParseUint(os.Args[3], 0, 64)
		Ck(err)
		write = true
	}

	mem, err := devmem.Open(target, size)
	Ck(err)

	res := mem.Read()
	Pf("read 0x%x\n", res)
	if write {
		mem.Write(newVal)
		Pf("write 0x%x\n", newVal)
		mem.Read()
	}

	mem.Close()

}
