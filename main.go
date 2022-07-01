package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"syscall"

	"github.com/dolmen-go/endian"
	. "github.com/stevegt/goadapt"
	"golang.org/x/sys/unix"
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

/*
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <signal.h>
#include <fcntl.h>
#include <ctype.h>
#include <termios.h>
#include <sys/types.h>
#include <sys/mman.h>

#define FATAL do { fprintf(stderr, "Error at line %d, file %s (%d) [%s]\n", \
  __LINE__, __FILE__, errno, strerror(errno)); exit(1); } while(0)


static inline void *fixup_addr(void *addr, size_t size);
*/

var order = endian.Native

var usage = `
Usage: %s { address } [ size [ data ] ]
	address : memory address to act upon
	size    : data bit width (default 32): 8, 16, 32, 64
	data    : data to be written
`

type Mem struct {
	page   []byte
	target int64
	size   int64
	sizeb  int
	base   int64
	diff   int
}

func MMap(fh *os.File, target, size int64) (m *Mem) {
	m = &Mem{}
	/* Map one page */
	pageSize := os.Getpagesize()
	m.base = target / int64(pageSize) * int64(pageSize)
	m.diff = int(target - m.base)
	Pf("target 0x%x pageSize 0x%x base 0x%x\n", target, pageSize, m.base)

	prot := syscall.PROT_READ | syscall.PROT_WRITE
	flags := syscall.MAP_SHARED
	fd := int(fh.Fd())

	var err error
	m.page, err = unix.Mmap(fd, m.base, pageSize, prot, flags)
	Ck(err)
	Pl("memory mapped")

	m.sizeb = int(size / 8)
	return
}

func (m *Mem) read() {
	raw := m.page[m.diff : m.diff+m.sizeb]
	Pf("raw 0x%x\n", raw)
	switch len(raw) * 8 {
	case 8:
		cooked := uint8(raw[0])
		// rd := bytes.NewReader(raw)
		// err := binary.Read(rd, order, &cooked)
		// Ck(err)
		Pf("read 0x%x\n", cooked)
	case 16:
		cooked := order.Uint16(raw)
		Pf("read 0x%x\n", cooked)
	case 32:
		cooked := order.Uint32(raw)
		Pf("read 0x%x\n", cooked)
	case 64:
		cooked := order.Uint64(raw)
		Pf("read 0x%x\n", cooked)
	default:
	}
}

func (m *Mem) write(newVal uint64) {
	wr := bytes.NewBuffer(make([]byte, 0, m.sizeb))
	err := binary.Write(wr, order, newVal)
	Ck(err)
	newRaw := wr.Bytes()
	for i := 0; i < m.sizeb; i++ {
		m.page[m.diff+i] = newRaw[i]
	}
	Pf("write 0x%x\n", newVal)
}

func (m *Mem) Unmap() {
	err := unix.Munmap(m.page)
	Ck(err)
}

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
		if !(size == 8 || size == 16 || size == 32 || size == 64) {
			Fpf(os.Stderr, "invalid bit width: %d\n", size)
			os.Exit(1)
		}
	}

	var newVal uint64
	write := false
	if len(os.Args) > 3 {
		newVal, err = strconv.ParseUint(os.Args[3], 0, 64)
		Ck(err)
		write = true
	}

	fh, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, os.ModePerm)
	Ck(err)
	defer fh.Close()
	Pl("/dev/mem opened")

	mem := MMap(fh, target, size)

	mem.read()
	if write {
		mem.write(newVal)
		mem.read()
	}

	mem.Unmap()

}
