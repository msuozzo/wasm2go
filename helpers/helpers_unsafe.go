//go:build ignore

package helpers

import (
	"encoding/binary"
	"math/bits"
	"unsafe"
)

// Faster memory access, using unsafe.
//
// The two-argument forms perform a single bounds check and then access the
// memory through the slice's data pointer. They require a non-negative
// address: for 32-bit memories the compiler computes addresses as
// uint32(index) [+ constant offset], so p is always in [0, 2³²+2³²) and
// p+size cannot overflow. The data pointer is re-derived from the slice
// header on every call, so it can never dangle across memory.grow.
//
// The slice forms (loadNNs/storeNNs) are used for 64-bit memories, whose
// wrap-around trap encoding can produce negative addresses that need the
// slice expression's own check.

//go:nosplit
func load16(m []byte, p int64) uint16 {
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
func store16(m []byte, p int64, v uint16) {
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
func load32(m []byte, p int64) uint32 {
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
func store32(m []byte, p int64, v uint32) {
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
func load64(m []byte, p int64) uint64 {
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
func store64(m []byte, p int64, v uint64) {
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
func load16s(b []byte) uint16 {
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
func store16s(b []byte, v uint16) {
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
func load32s(b []byte) uint32 {
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
func store32s(b []byte, v uint32) {
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
func load64s(b []byte) uint64 {
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
func store64s(b []byte, v uint64) {
	if !unalignedOK {
		binary.LittleEndian.PutUint64(b, v)
		return
	}
	if big {
		v = bits.ReverseBytes64(v)
	}
	*(*uint64)(unsafe.Pointer((*[8]byte)(b))) = v
}
