package devmem

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"

	"github.com/dolmen-go/endian"
	. "github.com/stevegt/goadapt"
	"golang.org/x/sys/unix"
)

var order = endian.Native

type Mem struct {
	page   []byte
	target int64
	size   int64
	sizeb  int
	base   int64
	diff   int
}

func Open(target, size int64) (m *Mem, err error) {
	defer Return(&err)
	m = &Mem{}

	fh, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, os.ModePerm)
	Ck(err)
	defer fh.Close()
	// Pl("/dev/mem opened")

	/* Map one page */
	pageSize := os.Getpagesize()
	m.base = target / int64(pageSize) * int64(pageSize)
	m.diff = int(target - m.base)
	// Pf("target 0x%x pageSize 0x%x base 0x%x\n", target, pageSize, m.base)

	prot := syscall.PROT_READ | syscall.PROT_WRITE
	flags := syscall.MAP_SHARED
	fd := int(fh.Fd())

	m.page, err = unix.Mmap(fd, m.base, pageSize, prot, flags)
	Ck(err)
	// Pl("memory mapped")

	if !(size == 8 || size == 16 || size == 32 || size == 64) {
		err = fmt.Errorf("invalid bit width: %d\n", size)
		return
	}

	m.size = size
	m.sizeb = int(size / 8)
	return
}

func (m *Mem) Read() (res uint64) {
	raw := m.page[m.diff : m.diff+m.sizeb]
	switch len(raw) * 8 {
	case 8:
		cooked := uint8(raw[0])
		res = uint64(cooked)
	case 16:
		cooked := order.Uint16(raw)
		res = uint64(cooked)
	case 32:
		cooked := order.Uint32(raw)
		res = uint64(cooked)
	case 64:
		cooked := order.Uint64(raw)
		res = uint64(cooked)
	default:
	}
	return
}

func (m *Mem) Write(newVal uint64) (err error) {
	defer Return(&err)
	wr := bytes.NewBuffer(make([]byte, 0, m.sizeb))
	err = binary.Write(wr, order, newVal)
	Ck(err)
	newRaw := wr.Bytes()
	for i := 0; i < m.sizeb; i++ {
		m.page[m.diff+i] = newRaw[i]
	}
	return
}

func (m *Mem) Close() (err error) {
	defer Return(&err)
	err = unix.Munmap(m.page)
	Ck(err)
	return
}
