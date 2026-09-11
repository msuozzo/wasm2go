//go:build ignore

package helpers

import (
	"encoding/binary"
	"math/bits"
	"unsafe"
)

// Faster memory access, using unsafe.

//go:nosplit
func load16[T uint32 | int64](mem []byte, addr T) uint16 {
	if !unalignedOK {
		return binary.LittleEndian.Uint16(mem[addr:])
	}
	_ = (*[2]byte)(mem[addr:])
	val := *(*uint16)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(mem)), uintptr(addr)))
	if big {
		return bits.ReverseBytes16(val)
	}
	return val
}

//go:nosplit
func store16[T uint32 | int64](mem []byte, addr T, val uint16) {
	if !unalignedOK {
		binary.LittleEndian.PutUint16(mem[addr:], val)
		return
	}
	if big {
		val = bits.ReverseBytes16(val)
	}
	_ = (*[2]byte)(mem[addr:])
	*(*uint16)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(mem)), uintptr(addr))) = val
}

//go:nosplit
func load32[T uint32 | int64](mem []byte, addr T) uint32 {
	if !unalignedOK {
		return binary.LittleEndian.Uint32(mem[addr:])
	}
	_ = (*[4]byte)(mem[addr:])
	val := *(*uint32)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(mem)), uintptr(addr)))
	if big {
		return bits.ReverseBytes32(val)
	}
	return val
}

//go:nosplit
func store32[T uint32 | int64](mem []byte, addr T, val uint32) {
	if !unalignedOK {
		binary.LittleEndian.PutUint32(mem[addr:], val)
		return
	}
	if big {
		val = bits.ReverseBytes32(val)
	}
	_ = (*[4]byte)(mem[addr:])
	*(*uint32)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(mem)), uintptr(addr))) = val
}

//go:nosplit
func load64[T uint32 | int64](mem []byte, addr T) uint64 {
	if !unalignedOK {
		return binary.LittleEndian.Uint64(mem[addr:])
	}
	_ = (*[8]byte)(mem[addr:])
	val := *(*uint64)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(mem)), uintptr(addr)))
	if big {
		return bits.ReverseBytes64(val)
	}
	return val
}

//go:nosplit
func store64[T uint32 | int64](mem []byte, addr T, val uint64) {
	if !unalignedOK {
		binary.LittleEndian.PutUint64(mem[addr:], val)
		return
	}
	if big {
		val = bits.ReverseBytes64(val)
	}
	_ = (*[8]byte)(mem[addr:])
	*(*uint64)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(mem)), uintptr(addr))) = val
}
