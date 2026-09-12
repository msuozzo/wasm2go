package libc

import "encoding/binary"

type ptr int32
type uptr uint32

var memory []byte

func load16(mem []byte, addr uptr) uint16 {
	return binary.LittleEndian.Uint16(mem[addr:])
}

func store16(mem []byte, addr uptr, val uint16) {
	binary.LittleEndian.PutUint16(mem[addr:], val)
}

func load32(mem []byte, addr uptr) uint32 {
	return binary.LittleEndian.Uint32(mem[addr:])
}

func store32(mem []byte, addr uptr, val uint32) {
	binary.LittleEndian.PutUint32(mem[addr:], val)
}

func load64(mem []byte, addr uptr) uint64 {
	return binary.LittleEndian.Uint64(mem[addr:])
}

func store64(mem []byte, addr uptr, val uint64) {
	binary.LittleEndian.PutUint64(mem[addr:], val)
}
