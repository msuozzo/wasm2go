//go:build ignore

package helpers

import (
	"encoding/binary"
	"math/bits"
	"unsafe"
)

// Faster memory access, using unsafe.
//
// loadNNat/storeNNat take the whole memory and an address p where p>0,
// p+size < INT64_MAX (the translator guarantees this for 32-bit, see memCall).
// They do one bounds check and then an unaligned access through the data
// pointer, re-derived from the slice header per call so it never dangles
// across memory.grow.
//
// loadNN/storeNN take a subslice, like the portable helpers, and support
// 64-bit and provided code that dynamically calls them (libc-gen output).

//go:nosplit
func load16at(m []byte, p int64) uint16 {
	_ = m[p+1]
	if !unalignedOK {
		return binary.LittleEndian.Uint16(m[p:])
	}
	v := *(*uint16)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(m)), uintptr(p)))
	if big {
		return bits.ReverseBytes16(v)
	}
	return v
}

//go:nosplit
func store16at(m []byte, p int64, v uint16) {
	_ = m[p+1]
	if !unalignedOK {
		binary.LittleEndian.PutUint16(m[p:], v)
		return
	}
	if big {
		v = bits.ReverseBytes16(v)
	}
	*(*uint16)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(m)), uintptr(p))) = v
}

//go:nosplit
func load32at(m []byte, p int64) uint32 {
	_ = m[p+3]
	if !unalignedOK {
		return binary.LittleEndian.Uint32(m[p:])
	}
	v := *(*uint32)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(m)), uintptr(p)))
	if big {
		return bits.ReverseBytes32(v)
	}
	return v
}

//go:nosplit
func store32at(m []byte, p int64, v uint32) {
	_ = m[p+3]
	if !unalignedOK {
		binary.LittleEndian.PutUint32(m[p:], v)
		return
	}
	if big {
		v = bits.ReverseBytes32(v)
	}
	*(*uint32)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(m)), uintptr(p))) = v
}

//go:nosplit
func load64at(m []byte, p int64) uint64 {
	_ = m[p+7]
	if !unalignedOK {
		return binary.LittleEndian.Uint64(m[p:])
	}
	v := *(*uint64)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(m)), uintptr(p)))
	if big {
		return bits.ReverseBytes64(v)
	}
	return v
}

//go:nosplit
func store64at(m []byte, p int64, v uint64) {
	_ = m[p+7]
	if !unalignedOK {
		binary.LittleEndian.PutUint64(m[p:], v)
		return
	}
	if big {
		v = bits.ReverseBytes64(v)
	}
	*(*uint64)(unsafe.Add(unsafe.Pointer(unsafe.SliceData(m)), uintptr(p))) = v
}

//go:nosplit
func load16(b []byte) uint16 {
	if !unalignedOK {
		return binary.LittleEndian.Uint16(b)
	}
	v := *(*uint16)(unsafe.Pointer((*[2]byte)(b)))
	if big {
		return bits.ReverseBytes16(v)
	}
	return v
}

//go:nosplit
func store16(b []byte, v uint16) {
	if !unalignedOK {
		binary.LittleEndian.PutUint16(b, v)
		return
	}
	if big {
		v = bits.ReverseBytes16(v)
	}
	*(*uint16)(unsafe.Pointer((*[2]byte)(b))) = v
}

//go:nosplit
func load32(b []byte) uint32 {
	if !unalignedOK {
		return binary.LittleEndian.Uint32(b)
	}
	v := *(*uint32)(unsafe.Pointer((*[4]byte)(b)))
	if big {
		return bits.ReverseBytes32(v)
	}
	return v
}

//go:nosplit
func store32(b []byte, v uint32) {
	if !unalignedOK {
		binary.LittleEndian.PutUint32(b, v)
		return
	}
	if big {
		v = bits.ReverseBytes32(v)
	}
	*(*uint32)(unsafe.Pointer((*[4]byte)(b))) = v
}

//go:nosplit
func load64(b []byte) uint64 {
	if !unalignedOK {
		return binary.LittleEndian.Uint64(b)
	}
	v := *(*uint64)(unsafe.Pointer((*[8]byte)(b)))
	if big {
		return bits.ReverseBytes64(v)
	}
	return v
}

//go:nosplit
func store64(b []byte, v uint64) {
	if !unalignedOK {
		binary.LittleEndian.PutUint64(b, v)
		return
	}
	if big {
		v = bits.ReverseBytes64(v)
	}
	*(*uint64)(unsafe.Pointer((*[8]byte)(b))) = v
}
