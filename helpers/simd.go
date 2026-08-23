package helpers

import (
	"encoding/binary"
	"math"
	"math/bits"
)

// The Wasm v128 type: 16 bytes, lanes stored little-endian,
// matching linear memory representation.
type v128 [16]byte

// Vector memory access.

//go:nosplit
func load128(b []byte) v128 {
	return v128(b)
}

//go:nosplit
func store128(b []byte, v v128) {
	copy(b[:16], v[:])
}

//go:nosplit
func v128_load8x8_s(x uint64) (r v128) {
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint16(r[2*i:], uint16(int16(int8(x>>(8*i)))))
	}
	return
}

//go:nosplit
func v128_load8x8_u(x uint64) (r v128) {
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint16(r[2*i:], uint16(byte(x>>(8*i))))
	}
	return
}

//go:nosplit
func v128_load16x4_s(x uint64) (r v128) {
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(r[4*i:], uint32(int32(int16(x>>(16*i)))))
	}
	return
}

//go:nosplit
func v128_load16x4_u(x uint64) (r v128) {
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(r[4*i:], uint32(uint16(x>>(16*i))))
	}
	return
}

//go:nosplit
func v128_load32x2_s(x uint64) (r v128) {
	for i := 0; i < 2; i++ {
		binary.LittleEndian.PutUint64(r[8*i:], uint64(int64(int32(x>>(32*i)))))
	}
	return
}

//go:nosplit
func v128_load32x2_u(x uint64) (r v128) {
	for i := 0; i < 2; i++ {
		binary.LittleEndian.PutUint64(r[8*i:], uint64(uint32(x>>(32*i))))
	}
	return
}

//go:nosplit
func v128_load32_zero(x uint32) (r v128) {
	binary.LittleEndian.PutUint32(r[0:], x)
	return
}

//go:nosplit
func v128_load64_zero(x uint64) (r v128) {
	binary.LittleEndian.PutUint64(r[0:], x)
	return
}

// Splats.

//go:nosplit
func i8x16_splat(x int32) (r v128) {
	for i := range r {
		r[i] = byte(x)
	}
	return
}

//go:nosplit
func i16x8_splat(x int32) (r v128) {
	for i := 0; i < 16; i += 2 {
		binary.LittleEndian.PutUint16(r[i:], uint16(x))
	}
	return
}

//go:nosplit
func i32x4_splat(x int32) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], uint32(x))
	}
	return
}

//go:nosplit
func i64x2_splat(x int64) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], uint64(x))
	}
	return
}

//go:nosplit
func f32x4_splat(x float32) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(x))
	}
	return
}

//go:nosplit
func f64x2_splat(x float64) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], math.Float64bits(x))
	}
	return
}

// Shuffles.

//go:nosplit
func i8x16_shuffle(a, b, m v128) (r v128) {
	for i := range r {
		if j := m[i]; j < 16 {
			r[i] = a[j]
		} else {
			r[i] = b[j-16]
		}
	}
	return
}

//go:nosplit
func i8x16_swizzle(a, s v128) (r v128) {
	for i := range r {
		if j := s[i]; j < 16 {
			r[i] = a[j]
		}
	}
	return
}

// Lane access.

//go:nosplit
func i8x16_extract_lane_s(v v128, l int) int32 {
	return int32(int8(v[l]))
}

//go:nosplit
func i8x16_extract_lane_u(v v128, l int) int32 {
	return int32(v[l])
}

//go:nosplit
func i8x16_replace_lane(v v128, l int, x int32) v128 {
	v[l] = byte(x)
	return v
}

//go:nosplit
func i16x8_extract_lane_s(v v128, l int) int32 {
	return int32(int16(binary.LittleEndian.Uint16(v[2*l:])))
}

//go:nosplit
func i16x8_extract_lane_u(v v128, l int) int32 {
	return int32(binary.LittleEndian.Uint16(v[2*l:]))
}

//go:nosplit
func i16x8_replace_lane(v v128, l int, x int32) v128 {
	binary.LittleEndian.PutUint16(v[2*l:], uint16(x))
	return v
}

//go:nosplit
func i32x4_extract_lane(v v128, l int) int32 {
	return int32(binary.LittleEndian.Uint32(v[4*l:]))
}

//go:nosplit
func i32x4_replace_lane(v v128, l int, x int32) v128 {
	binary.LittleEndian.PutUint32(v[4*l:], uint32(x))
	return v
}

//go:nosplit
func i64x2_extract_lane(v v128, l int) int64 {
	return int64(binary.LittleEndian.Uint64(v[8*l:]))
}

//go:nosplit
func i64x2_replace_lane(v v128, l int, x int64) v128 {
	binary.LittleEndian.PutUint64(v[8*l:], uint64(x))
	return v
}

//go:nosplit
func f32x4_extract_lane(v v128, l int) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(v[4*l:]))
}

//go:nosplit
func f32x4_replace_lane(v v128, l int, x float32) v128 {
	binary.LittleEndian.PutUint32(v[4*l:], math.Float32bits(x))
	return v
}

//go:nosplit
func f64x2_extract_lane(v v128, l int) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(v[8*l:]))
}

//go:nosplit
func f64x2_replace_lane(v v128, l int, x float64) v128 {
	binary.LittleEndian.PutUint64(v[8*l:], math.Float64bits(x))
	return v
}

// Bitwise operations.

//go:nosplit
func v128_not(a v128) (r v128) {
	for i := range r {
		r[i] = ^a[i]
	}
	return
}

//go:nosplit
func v128_and(a, b v128) (r v128) {
	for i := range r {
		r[i] = a[i] & b[i]
	}
	return
}

//go:nosplit
func v128_andnot(a, b v128) (r v128) {
	for i := range r {
		r[i] = a[i] &^ b[i]
	}
	return
}

//go:nosplit
func v128_or(a, b v128) (r v128) {
	for i := range r {
		r[i] = a[i] | b[i]
	}
	return
}

//go:nosplit
func v128_xor(a, b v128) (r v128) {
	for i := range r {
		r[i] = a[i] ^ b[i]
	}
	return
}

//go:nosplit
func v128_bitselect(a, b, c v128) (r v128) {
	for i := range r {
		r[i] = a[i]&c[i] | b[i]&^c[i]
	}
	return
}

//go:nosplit
func v128_any_true(v v128) int32 {
	var zero v128
	if v != zero {
		return 1
	}
	return 0
}

// Boolean reductions and bitmasks.

//go:nosplit
func i8x16_all_true(v v128) int32 {
	for i := range v {
		if v[i] == 0 {
			return 0
		}
	}
	return 1
}

//go:nosplit
func i16x8_all_true(v v128) int32 {
	for i := 0; i < 16; i += 2 {
		if binary.LittleEndian.Uint16(v[i:]) == 0 {
			return 0
		}
	}
	return 1
}

//go:nosplit
func i32x4_all_true(v v128) int32 {
	for i := 0; i < 16; i += 4 {
		if binary.LittleEndian.Uint32(v[i:]) == 0 {
			return 0
		}
	}
	return 1
}

//go:nosplit
func i64x2_all_true(v v128) int32 {
	for i := 0; i < 16; i += 8 {
		if binary.LittleEndian.Uint64(v[i:]) == 0 {
			return 0
		}
	}
	return 1
}

//go:nosplit
func i8x16_bitmask(v v128) (r int32) {
	for i := range v {
		r |= int32(v[i]>>7) << i
	}
	return
}

//go:nosplit
func i16x8_bitmask(v v128) (r int32) {
	for i := 0; i < 8; i++ {
		r |= int32(binary.LittleEndian.Uint16(v[2*i:])>>15) << i
	}
	return
}

//go:nosplit
func i32x4_bitmask(v v128) (r int32) {
	for i := 0; i < 4; i++ {
		r |= int32(binary.LittleEndian.Uint32(v[4*i:])>>31) << i
	}
	return
}

//go:nosplit
func i64x2_bitmask(v v128) (r int32) {
	for i := 0; i < 2; i++ {
		r |= int32(binary.LittleEndian.Uint64(v[8*i:])>>63) << i
	}
	return
}

// Integer comparisons.

//go:nosplit
func i8x16_eq(a, b v128) (r v128) {
	for i := range r {
		if a[i] == b[i] {
			r[i] = 0xff
		}
	}
	return
}

//go:nosplit
func i8x16_ne(a, b v128) (r v128) {
	for i := range r {
		if a[i] != b[i] {
			r[i] = 0xff
		}
	}
	return
}

//go:nosplit
func i8x16_lt_s(a, b v128) (r v128) {
	for i := range r {
		if int8(a[i]) < int8(b[i]) {
			r[i] = 0xff
		}
	}
	return
}

//go:nosplit
func i8x16_lt_u(a, b v128) (r v128) {
	for i := range r {
		if a[i] < b[i] {
			r[i] = 0xff
		}
	}
	return
}

//go:nosplit
func i8x16_gt_s(a, b v128) (r v128) {
	for i := range r {
		if int8(a[i]) > int8(b[i]) {
			r[i] = 0xff
		}
	}
	return
}

//go:nosplit
func i8x16_gt_u(a, b v128) (r v128) {
	for i := range r {
		if a[i] > b[i] {
			r[i] = 0xff
		}
	}
	return
}

//go:nosplit
func i8x16_le_s(a, b v128) (r v128) {
	for i := range r {
		if int8(a[i]) <= int8(b[i]) {
			r[i] = 0xff
		}
	}
	return
}

//go:nosplit
func i8x16_le_u(a, b v128) (r v128) {
	for i := range r {
		if a[i] <= b[i] {
			r[i] = 0xff
		}
	}
	return
}

//go:nosplit
func i8x16_ge_s(a, b v128) (r v128) {
	for i := range r {
		if int8(a[i]) >= int8(b[i]) {
			r[i] = 0xff
		}
	}
	return
}

//go:nosplit
func i8x16_ge_u(a, b v128) (r v128) {
	for i := range r {
		if a[i] >= b[i] {
			r[i] = 0xff
		}
	}
	return
}

//go:nosplit
func i16x8_eq(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		if binary.LittleEndian.Uint16(a[i:]) == binary.LittleEndian.Uint16(b[i:]) {
			binary.LittleEndian.PutUint16(r[i:], 0xffff)
		}
	}
	return
}

//go:nosplit
func i16x8_ne(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		if binary.LittleEndian.Uint16(a[i:]) != binary.LittleEndian.Uint16(b[i:]) {
			binary.LittleEndian.PutUint16(r[i:], 0xffff)
		}
	}
	return
}

//go:nosplit
func i16x8_lt_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		if int16(binary.LittleEndian.Uint16(a[i:])) < int16(binary.LittleEndian.Uint16(b[i:])) {
			binary.LittleEndian.PutUint16(r[i:], 0xffff)
		}
	}
	return
}

//go:nosplit
func i16x8_lt_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		if binary.LittleEndian.Uint16(a[i:]) < binary.LittleEndian.Uint16(b[i:]) {
			binary.LittleEndian.PutUint16(r[i:], 0xffff)
		}
	}
	return
}

//go:nosplit
func i16x8_gt_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		if int16(binary.LittleEndian.Uint16(a[i:])) > int16(binary.LittleEndian.Uint16(b[i:])) {
			binary.LittleEndian.PutUint16(r[i:], 0xffff)
		}
	}
	return
}

//go:nosplit
func i16x8_gt_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		if binary.LittleEndian.Uint16(a[i:]) > binary.LittleEndian.Uint16(b[i:]) {
			binary.LittleEndian.PutUint16(r[i:], 0xffff)
		}
	}
	return
}

//go:nosplit
func i16x8_le_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		if int16(binary.LittleEndian.Uint16(a[i:])) <= int16(binary.LittleEndian.Uint16(b[i:])) {
			binary.LittleEndian.PutUint16(r[i:], 0xffff)
		}
	}
	return
}

//go:nosplit
func i16x8_le_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		if binary.LittleEndian.Uint16(a[i:]) <= binary.LittleEndian.Uint16(b[i:]) {
			binary.LittleEndian.PutUint16(r[i:], 0xffff)
		}
	}
	return
}

//go:nosplit
func i16x8_ge_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		if int16(binary.LittleEndian.Uint16(a[i:])) >= int16(binary.LittleEndian.Uint16(b[i:])) {
			binary.LittleEndian.PutUint16(r[i:], 0xffff)
		}
	}
	return
}

//go:nosplit
func i16x8_ge_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		if binary.LittleEndian.Uint16(a[i:]) >= binary.LittleEndian.Uint16(b[i:]) {
			binary.LittleEndian.PutUint16(r[i:], 0xffff)
		}
	}
	return
}

//go:nosplit
func i32x4_eq(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if binary.LittleEndian.Uint32(a[i:]) == binary.LittleEndian.Uint32(b[i:]) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func i32x4_ne(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if binary.LittleEndian.Uint32(a[i:]) != binary.LittleEndian.Uint32(b[i:]) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func i32x4_lt_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if int32(binary.LittleEndian.Uint32(a[i:])) < int32(binary.LittleEndian.Uint32(b[i:])) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func i32x4_lt_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if binary.LittleEndian.Uint32(a[i:]) < binary.LittleEndian.Uint32(b[i:]) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func i32x4_gt_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if int32(binary.LittleEndian.Uint32(a[i:])) > int32(binary.LittleEndian.Uint32(b[i:])) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func i32x4_gt_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if binary.LittleEndian.Uint32(a[i:]) > binary.LittleEndian.Uint32(b[i:]) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func i32x4_le_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if int32(binary.LittleEndian.Uint32(a[i:])) <= int32(binary.LittleEndian.Uint32(b[i:])) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func i32x4_le_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if binary.LittleEndian.Uint32(a[i:]) <= binary.LittleEndian.Uint32(b[i:]) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func i32x4_ge_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if int32(binary.LittleEndian.Uint32(a[i:])) >= int32(binary.LittleEndian.Uint32(b[i:])) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func i32x4_ge_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if binary.LittleEndian.Uint32(a[i:]) >= binary.LittleEndian.Uint32(b[i:]) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func i64x2_eq(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if binary.LittleEndian.Uint64(a[i:]) == binary.LittleEndian.Uint64(b[i:]) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

//go:nosplit
func i64x2_ne(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if binary.LittleEndian.Uint64(a[i:]) != binary.LittleEndian.Uint64(b[i:]) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

//go:nosplit
func i64x2_lt_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if int64(binary.LittleEndian.Uint64(a[i:])) < int64(binary.LittleEndian.Uint64(b[i:])) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

//go:nosplit
func i64x2_gt_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if int64(binary.LittleEndian.Uint64(a[i:])) > int64(binary.LittleEndian.Uint64(b[i:])) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

//go:nosplit
func i64x2_le_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if int64(binary.LittleEndian.Uint64(a[i:])) <= int64(binary.LittleEndian.Uint64(b[i:])) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

//go:nosplit
func i64x2_ge_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if int64(binary.LittleEndian.Uint64(a[i:])) >= int64(binary.LittleEndian.Uint64(b[i:])) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

// Float comparisons.

//go:nosplit
func f32x4_eq(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if math.Float32frombits(binary.LittleEndian.Uint32(a[i:])) == math.Float32frombits(binary.LittleEndian.Uint32(b[i:])) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func f32x4_ne(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if math.Float32frombits(binary.LittleEndian.Uint32(a[i:])) != math.Float32frombits(binary.LittleEndian.Uint32(b[i:])) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func f32x4_lt(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if math.Float32frombits(binary.LittleEndian.Uint32(a[i:])) < math.Float32frombits(binary.LittleEndian.Uint32(b[i:])) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func f32x4_gt(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if math.Float32frombits(binary.LittleEndian.Uint32(a[i:])) > math.Float32frombits(binary.LittleEndian.Uint32(b[i:])) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func f32x4_le(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if math.Float32frombits(binary.LittleEndian.Uint32(a[i:])) <= math.Float32frombits(binary.LittleEndian.Uint32(b[i:])) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func f32x4_ge(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		if math.Float32frombits(binary.LittleEndian.Uint32(a[i:])) >= math.Float32frombits(binary.LittleEndian.Uint32(b[i:])) {
			binary.LittleEndian.PutUint32(r[i:], 0xffffffff)
		}
	}
	return
}

//go:nosplit
func f64x2_eq(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if math.Float64frombits(binary.LittleEndian.Uint64(a[i:])) == math.Float64frombits(binary.LittleEndian.Uint64(b[i:])) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

//go:nosplit
func f64x2_ne(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if math.Float64frombits(binary.LittleEndian.Uint64(a[i:])) != math.Float64frombits(binary.LittleEndian.Uint64(b[i:])) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

//go:nosplit
func f64x2_lt(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if math.Float64frombits(binary.LittleEndian.Uint64(a[i:])) < math.Float64frombits(binary.LittleEndian.Uint64(b[i:])) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

//go:nosplit
func f64x2_gt(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if math.Float64frombits(binary.LittleEndian.Uint64(a[i:])) > math.Float64frombits(binary.LittleEndian.Uint64(b[i:])) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

//go:nosplit
func f64x2_le(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if math.Float64frombits(binary.LittleEndian.Uint64(a[i:])) <= math.Float64frombits(binary.LittleEndian.Uint64(b[i:])) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

//go:nosplit
func f64x2_ge(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		if math.Float64frombits(binary.LittleEndian.Uint64(a[i:])) >= math.Float64frombits(binary.LittleEndian.Uint64(b[i:])) {
			binary.LittleEndian.PutUint64(r[i:], 0xffffffffffffffff)
		}
	}
	return
}

// Integer unary operations.

//go:nosplit
func i8x16_abs(a v128) (r v128) {
	for i := range r {
		x := int8(a[i])
		if x < 0 {
			x = -x
		}
		r[i] = byte(x)
	}
	return
}

//go:nosplit
func i8x16_neg(a v128) (r v128) {
	for i := range r {
		r[i] = byte(-int8(a[i]))
	}
	return
}

//go:nosplit
func i8x16_popcnt(a v128) (r v128) {
	for i := range r {
		r[i] = byte(bits.OnesCount8(a[i]))
	}
	return
}

//go:nosplit
func i16x8_abs(a v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := int16(binary.LittleEndian.Uint16(a[i:]))
		if x < 0 {
			x = -x
		}
		binary.LittleEndian.PutUint16(r[i:], uint16(x))
	}
	return
}

//go:nosplit
func i16x8_neg(a v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		binary.LittleEndian.PutUint16(r[i:], uint16(-int16(binary.LittleEndian.Uint16(a[i:]))))
	}
	return
}

//go:nosplit
func i32x4_abs(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := int32(binary.LittleEndian.Uint32(a[i:]))
		if x < 0 {
			x = -x
		}
		binary.LittleEndian.PutUint32(r[i:], uint32(x))
	}
	return
}

//go:nosplit
func i32x4_neg(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], uint32(-int32(binary.LittleEndian.Uint32(a[i:]))))
	}
	return
}

//go:nosplit
func i64x2_abs(a v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := int64(binary.LittleEndian.Uint64(a[i:]))
		if x < 0 {
			x = -x
		}
		binary.LittleEndian.PutUint64(r[i:], uint64(x))
	}
	return
}

//go:nosplit
func i64x2_neg(a v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], uint64(-int64(binary.LittleEndian.Uint64(a[i:]))))
	}
	return
}

// Integer shifts: the count is taken modulo lane width.

//go:nosplit
func i8x16_shl(a v128, y int32) (r v128) {
	for i := range r {
		r[i] = a[i] << (y & 7)
	}
	return
}

//go:nosplit
func i8x16_shr_s(a v128, y int32) (r v128) {
	for i := range r {
		r[i] = byte(int8(a[i]) >> (y & 7))
	}
	return
}

//go:nosplit
func i8x16_shr_u(a v128, y int32) (r v128) {
	for i := range r {
		r[i] = a[i] >> (y & 7)
	}
	return
}

//go:nosplit
func i16x8_shl(a v128, y int32) (r v128) {
	for i := 0; i < 16; i += 2 {
		binary.LittleEndian.PutUint16(r[i:], binary.LittleEndian.Uint16(a[i:])<<(y&15))
	}
	return
}

//go:nosplit
func i16x8_shr_s(a v128, y int32) (r v128) {
	for i := 0; i < 16; i += 2 {
		binary.LittleEndian.PutUint16(r[i:], uint16(int16(binary.LittleEndian.Uint16(a[i:]))>>(y&15)))
	}
	return
}

//go:nosplit
func i16x8_shr_u(a v128, y int32) (r v128) {
	for i := 0; i < 16; i += 2 {
		binary.LittleEndian.PutUint16(r[i:], binary.LittleEndian.Uint16(a[i:])>>(y&15))
	}
	return
}

//go:nosplit
func i32x4_shl(a v128, y int32) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], binary.LittleEndian.Uint32(a[i:])<<(y&31))
	}
	return
}

//go:nosplit
func i32x4_shr_s(a v128, y int32) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], uint32(int32(binary.LittleEndian.Uint32(a[i:]))>>(y&31)))
	}
	return
}

//go:nosplit
func i32x4_shr_u(a v128, y int32) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], binary.LittleEndian.Uint32(a[i:])>>(y&31))
	}
	return
}

//go:nosplit
func i64x2_shl(a v128, y int32) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], binary.LittleEndian.Uint64(a[i:])<<(y&63))
	}
	return
}

//go:nosplit
func i64x2_shr_s(a v128, y int32) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], uint64(int64(binary.LittleEndian.Uint64(a[i:]))>>(y&63)))
	}
	return
}

//go:nosplit
func i64x2_shr_u(a v128, y int32) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], binary.LittleEndian.Uint64(a[i:])>>(y&63))
	}
	return
}

// Integer arithmetic.

//go:nosplit
func i8x16_add(a, b v128) (r v128) {
	for i := range r {
		r[i] = a[i] + b[i]
	}
	return
}

//go:nosplit
func i8x16_sub(a, b v128) (r v128) {
	for i := range r {
		r[i] = a[i] - b[i]
	}
	return
}

//go:nosplit
func i8x16_add_sat_s(a, b v128) (r v128) {
	for i := range r {
		x := int32(int8(a[i])) + int32(int8(b[i]))
		if x < math.MinInt8 {
			x = math.MinInt8
		}
		if x > math.MaxInt8 {
			x = math.MaxInt8
		}
		r[i] = byte(x)
	}
	return
}

//go:nosplit
func i8x16_add_sat_u(a, b v128) (r v128) {
	for i := range r {
		x := int32(a[i]) + int32(b[i])
		if x > math.MaxUint8 {
			x = math.MaxUint8
		}
		r[i] = byte(x)
	}
	return
}

//go:nosplit
func i8x16_sub_sat_s(a, b v128) (r v128) {
	for i := range r {
		x := int32(int8(a[i])) - int32(int8(b[i]))
		if x < math.MinInt8 {
			x = math.MinInt8
		}
		if x > math.MaxInt8 {
			x = math.MaxInt8
		}
		r[i] = byte(x)
	}
	return
}

//go:nosplit
func i8x16_sub_sat_u(a, b v128) (r v128) {
	for i := range r {
		if a[i] > b[i] {
			r[i] = a[i] - b[i]
		}
	}
	return
}

//go:nosplit
func i8x16_min_s(a, b v128) (r v128) {
	for i := range r {
		if int8(a[i]) < int8(b[i]) {
			r[i] = a[i]
		} else {
			r[i] = b[i]
		}
	}
	return
}

//go:nosplit
func i8x16_min_u(a, b v128) (r v128) {
	for i := range r {
		r[i] = min(a[i], b[i])
	}
	return
}

//go:nosplit
func i8x16_max_s(a, b v128) (r v128) {
	for i := range r {
		if int8(a[i]) > int8(b[i]) {
			r[i] = a[i]
		} else {
			r[i] = b[i]
		}
	}
	return
}

//go:nosplit
func i8x16_max_u(a, b v128) (r v128) {
	for i := range r {
		r[i] = max(a[i], b[i])
	}
	return
}

//go:nosplit
func i8x16_avgr_u(a, b v128) (r v128) {
	for i := range r {
		r[i] = byte((int32(a[i]) + int32(b[i]) + 1) / 2)
	}
	return
}

//go:nosplit
func i16x8_add(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		binary.LittleEndian.PutUint16(r[i:], binary.LittleEndian.Uint16(a[i:])+binary.LittleEndian.Uint16(b[i:]))
	}
	return
}

//go:nosplit
func i16x8_sub(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		binary.LittleEndian.PutUint16(r[i:], binary.LittleEndian.Uint16(a[i:])-binary.LittleEndian.Uint16(b[i:]))
	}
	return
}

//go:nosplit
func i16x8_mul(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		binary.LittleEndian.PutUint16(r[i:], binary.LittleEndian.Uint16(a[i:])*binary.LittleEndian.Uint16(b[i:]))
	}
	return
}

//go:nosplit
func i16x8_add_sat_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := int32(int16(binary.LittleEndian.Uint16(a[i:]))) + int32(int16(binary.LittleEndian.Uint16(b[i:])))
		if x < math.MinInt16 {
			x = math.MinInt16
		}
		if x > math.MaxInt16 {
			x = math.MaxInt16
		}
		binary.LittleEndian.PutUint16(r[i:], uint16(x))
	}
	return
}

//go:nosplit
func i16x8_add_sat_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := int32(binary.LittleEndian.Uint16(a[i:])) + int32(binary.LittleEndian.Uint16(b[i:]))
		if x > math.MaxUint16 {
			x = math.MaxUint16
		}
		binary.LittleEndian.PutUint16(r[i:], uint16(x))
	}
	return
}

//go:nosplit
func i16x8_sub_sat_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := int32(int16(binary.LittleEndian.Uint16(a[i:]))) - int32(int16(binary.LittleEndian.Uint16(b[i:])))
		if x < math.MinInt16 {
			x = math.MinInt16
		}
		if x > math.MaxInt16 {
			x = math.MaxInt16
		}
		binary.LittleEndian.PutUint16(r[i:], uint16(x))
	}
	return
}

//go:nosplit
func i16x8_sub_sat_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := binary.LittleEndian.Uint16(a[i:])
		y := binary.LittleEndian.Uint16(b[i:])
		if x > y {
			binary.LittleEndian.PutUint16(r[i:], x-y)
		}
	}
	return
}

//go:nosplit
func i16x8_min_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := int16(binary.LittleEndian.Uint16(a[i:]))
		y := int16(binary.LittleEndian.Uint16(b[i:]))
		binary.LittleEndian.PutUint16(r[i:], uint16(min(x, y)))
	}
	return
}

//go:nosplit
func i16x8_min_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		binary.LittleEndian.PutUint16(r[i:], min(binary.LittleEndian.Uint16(a[i:]), binary.LittleEndian.Uint16(b[i:])))
	}
	return
}

//go:nosplit
func i16x8_max_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := int16(binary.LittleEndian.Uint16(a[i:]))
		y := int16(binary.LittleEndian.Uint16(b[i:]))
		binary.LittleEndian.PutUint16(r[i:], uint16(max(x, y)))
	}
	return
}

//go:nosplit
func i16x8_max_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		binary.LittleEndian.PutUint16(r[i:], max(binary.LittleEndian.Uint16(a[i:]), binary.LittleEndian.Uint16(b[i:])))
	}
	return
}

//go:nosplit
func i16x8_avgr_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := int32(binary.LittleEndian.Uint16(a[i:])) + int32(binary.LittleEndian.Uint16(b[i:])) + 1
		binary.LittleEndian.PutUint16(r[i:], uint16(x/2))
	}
	return
}

//go:nosplit
func i16x8_q15mulr_sat_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := int32(int16(binary.LittleEndian.Uint16(a[i:])))
		y := int32(int16(binary.LittleEndian.Uint16(b[i:])))
		p := (x*y + 0x4000) >> 15
		if p > math.MaxInt16 {
			p = math.MaxInt16
		}
		binary.LittleEndian.PutUint16(r[i:], uint16(p))
	}
	return
}

//go:nosplit
func i32x4_add(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], binary.LittleEndian.Uint32(a[i:])+binary.LittleEndian.Uint32(b[i:]))
	}
	return
}

//go:nosplit
func i32x4_sub(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], binary.LittleEndian.Uint32(a[i:])-binary.LittleEndian.Uint32(b[i:]))
	}
	return
}

//go:nosplit
func i32x4_mul(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], binary.LittleEndian.Uint32(a[i:])*binary.LittleEndian.Uint32(b[i:]))
	}
	return
}

//go:nosplit
func i32x4_min_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := int32(binary.LittleEndian.Uint32(a[i:]))
		y := int32(binary.LittleEndian.Uint32(b[i:]))
		binary.LittleEndian.PutUint32(r[i:], uint32(min(x, y)))
	}
	return
}

//go:nosplit
func i32x4_min_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], min(binary.LittleEndian.Uint32(a[i:]), binary.LittleEndian.Uint32(b[i:])))
	}
	return
}

//go:nosplit
func i32x4_max_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := int32(binary.LittleEndian.Uint32(a[i:]))
		y := int32(binary.LittleEndian.Uint32(b[i:]))
		binary.LittleEndian.PutUint32(r[i:], uint32(max(x, y)))
	}
	return
}

//go:nosplit
func i32x4_max_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], max(binary.LittleEndian.Uint32(a[i:]), binary.LittleEndian.Uint32(b[i:])))
	}
	return
}

//go:nosplit
func i32x4_dot_i16x8_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x0 := int32(int16(binary.LittleEndian.Uint16(a[i:])))
		y0 := int32(int16(binary.LittleEndian.Uint16(b[i:])))
		x1 := int32(int16(binary.LittleEndian.Uint16(a[i+2:])))
		y1 := int32(int16(binary.LittleEndian.Uint16(b[i+2:])))
		binary.LittleEndian.PutUint32(r[i:], uint32(x0*y0+x1*y1))
	}
	return
}

//go:nosplit
func i64x2_add(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], binary.LittleEndian.Uint64(a[i:])+binary.LittleEndian.Uint64(b[i:]))
	}
	return
}

//go:nosplit
func i64x2_sub(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], binary.LittleEndian.Uint64(a[i:])-binary.LittleEndian.Uint64(b[i:]))
	}
	return
}

//go:nosplit
func i64x2_mul(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], binary.LittleEndian.Uint64(a[i:])*binary.LittleEndian.Uint64(b[i:]))
	}
	return
}

// Narrowing conversions: saturate lanes of a, then b, into a narrower result.

//go:nosplit
func i8x16_narrow_i16x8_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := int16(binary.LittleEndian.Uint16(a[i:]))
		if x < math.MinInt8 {
			x = math.MinInt8
		}
		if x > math.MaxInt8 {
			x = math.MaxInt8
		}
		r[i/2] = byte(x)
		x = int16(binary.LittleEndian.Uint16(b[i:]))
		if x < math.MinInt8 {
			x = math.MinInt8
		}
		if x > math.MaxInt8 {
			x = math.MaxInt8
		}
		r[8+i/2] = byte(x)
	}
	return
}

//go:nosplit
func i8x16_narrow_i16x8_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 2 {
		x := int16(binary.LittleEndian.Uint16(a[i:]))
		if x < 0 {
			x = 0
		}
		if x > math.MaxUint8 {
			x = math.MaxUint8
		}
		r[i/2] = byte(x)
		x = int16(binary.LittleEndian.Uint16(b[i:]))
		if x < 0 {
			x = 0
		}
		if x > math.MaxUint8 {
			x = math.MaxUint8
		}
		r[8+i/2] = byte(x)
	}
	return
}

//go:nosplit
func i16x8_narrow_i32x4_s(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := int32(binary.LittleEndian.Uint32(a[i:]))
		if x < math.MinInt16 {
			x = math.MinInt16
		}
		if x > math.MaxInt16 {
			x = math.MaxInt16
		}
		binary.LittleEndian.PutUint16(r[i/2:], uint16(x))
		x = int32(binary.LittleEndian.Uint32(b[i:]))
		if x < math.MinInt16 {
			x = math.MinInt16
		}
		if x > math.MaxInt16 {
			x = math.MaxInt16
		}
		binary.LittleEndian.PutUint16(r[8+i/2:], uint16(x))
	}
	return
}

//go:nosplit
func i16x8_narrow_i32x4_u(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := int32(binary.LittleEndian.Uint32(a[i:]))
		if x < 0 {
			x = 0
		}
		if x > math.MaxUint16 {
			x = math.MaxUint16
		}
		binary.LittleEndian.PutUint16(r[i/2:], uint16(x))
		x = int32(binary.LittleEndian.Uint32(b[i:]))
		if x < 0 {
			x = 0
		}
		if x > math.MaxUint16 {
			x = math.MaxUint16
		}
		binary.LittleEndian.PutUint16(r[8+i/2:], uint16(x))
	}
	return
}

// Widening conversions.

//go:nosplit
func i16x8_extend_low_i8x16_s(a v128) (r v128) {
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint16(r[2*i:], uint16(int16(int8(a[i]))))
	}
	return
}

//go:nosplit
func i16x8_extend_high_i8x16_s(a v128) (r v128) {
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint16(r[2*i:], uint16(int16(int8(a[8+i]))))
	}
	return
}

//go:nosplit
func i16x8_extend_low_i8x16_u(a v128) (r v128) {
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint16(r[2*i:], uint16(a[i]))
	}
	return
}

//go:nosplit
func i16x8_extend_high_i8x16_u(a v128) (r v128) {
	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint16(r[2*i:], uint16(a[8+i]))
	}
	return
}

//go:nosplit
func i32x4_extend_low_i16x8_s(a v128) (r v128) {
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(r[4*i:], uint32(int32(int16(binary.LittleEndian.Uint16(a[2*i:])))))
	}
	return
}

//go:nosplit
func i32x4_extend_high_i16x8_s(a v128) (r v128) {
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(r[4*i:], uint32(int32(int16(binary.LittleEndian.Uint16(a[8+2*i:])))))
	}
	return
}

//go:nosplit
func i32x4_extend_low_i16x8_u(a v128) (r v128) {
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(r[4*i:], uint32(binary.LittleEndian.Uint16(a[2*i:])))
	}
	return
}

//go:nosplit
func i32x4_extend_high_i16x8_u(a v128) (r v128) {
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(r[4*i:], uint32(binary.LittleEndian.Uint16(a[8+2*i:])))
	}
	return
}

//go:nosplit
func i64x2_extend_low_i32x4_s(a v128) (r v128) {
	for i := 0; i < 2; i++ {
		binary.LittleEndian.PutUint64(r[8*i:], uint64(int64(int32(binary.LittleEndian.Uint32(a[4*i:])))))
	}
	return
}

//go:nosplit
func i64x2_extend_high_i32x4_s(a v128) (r v128) {
	for i := 0; i < 2; i++ {
		binary.LittleEndian.PutUint64(r[8*i:], uint64(int64(int32(binary.LittleEndian.Uint32(a[8+4*i:])))))
	}
	return
}

//go:nosplit
func i64x2_extend_low_i32x4_u(a v128) (r v128) {
	for i := 0; i < 2; i++ {
		binary.LittleEndian.PutUint64(r[8*i:], uint64(binary.LittleEndian.Uint32(a[4*i:])))
	}
	return
}

//go:nosplit
func i64x2_extend_high_i32x4_u(a v128) (r v128) {
	for i := 0; i < 2; i++ {
		binary.LittleEndian.PutUint64(r[8*i:], uint64(binary.LittleEndian.Uint32(a[8+4*i:])))
	}
	return
}

// Extended pairwise additions.

//go:nosplit
func i16x8_extadd_pairwise_i8x16_s(a v128) (r v128) {
	for i := 0; i < 8; i++ {
		x := int16(int8(a[2*i])) + int16(int8(a[2*i+1]))
		binary.LittleEndian.PutUint16(r[2*i:], uint16(x))
	}
	return
}

//go:nosplit
func i16x8_extadd_pairwise_i8x16_u(a v128) (r v128) {
	for i := 0; i < 8; i++ {
		x := uint16(a[2*i]) + uint16(a[2*i+1])
		binary.LittleEndian.PutUint16(r[2*i:], x)
	}
	return
}

//go:nosplit
func i32x4_extadd_pairwise_i16x8_s(a v128) (r v128) {
	for i := 0; i < 4; i++ {
		x := int32(int16(binary.LittleEndian.Uint16(a[4*i:]))) + int32(int16(binary.LittleEndian.Uint16(a[4*i+2:])))
		binary.LittleEndian.PutUint32(r[4*i:], uint32(x))
	}
	return
}

//go:nosplit
func i32x4_extadd_pairwise_i16x8_u(a v128) (r v128) {
	for i := 0; i < 4; i++ {
		x := uint32(binary.LittleEndian.Uint16(a[4*i:])) + uint32(binary.LittleEndian.Uint16(a[4*i+2:]))
		binary.LittleEndian.PutUint32(r[4*i:], x)
	}
	return
}

// Extended multiplications.

//go:nosplit
func i16x8_extmul_low_i8x16_s(a, b v128) (r v128) {
	for i := 0; i < 8; i++ {
		x := int16(int8(a[i])) * int16(int8(b[i]))
		binary.LittleEndian.PutUint16(r[2*i:], uint16(x))
	}
	return
}

//go:nosplit
func i16x8_extmul_high_i8x16_s(a, b v128) (r v128) {
	for i := 0; i < 8; i++ {
		x := int16(int8(a[8+i])) * int16(int8(b[8+i]))
		binary.LittleEndian.PutUint16(r[2*i:], uint16(x))
	}
	return
}

//go:nosplit
func i16x8_extmul_low_i8x16_u(a, b v128) (r v128) {
	for i := 0; i < 8; i++ {
		x := uint16(a[i]) * uint16(b[i])
		binary.LittleEndian.PutUint16(r[2*i:], x)
	}
	return
}

//go:nosplit
func i16x8_extmul_high_i8x16_u(a, b v128) (r v128) {
	for i := 0; i < 8; i++ {
		x := uint16(a[8+i]) * uint16(b[8+i])
		binary.LittleEndian.PutUint16(r[2*i:], x)
	}
	return
}

//go:nosplit
func i32x4_extmul_low_i16x8_s(a, b v128) (r v128) {
	for i := 0; i < 4; i++ {
		x := int32(int16(binary.LittleEndian.Uint16(a[2*i:]))) * int32(int16(binary.LittleEndian.Uint16(b[2*i:])))
		binary.LittleEndian.PutUint32(r[4*i:], uint32(x))
	}
	return
}

//go:nosplit
func i32x4_extmul_high_i16x8_s(a, b v128) (r v128) {
	for i := 0; i < 4; i++ {
		x := int32(int16(binary.LittleEndian.Uint16(a[8+2*i:]))) * int32(int16(binary.LittleEndian.Uint16(b[8+2*i:])))
		binary.LittleEndian.PutUint32(r[4*i:], uint32(x))
	}
	return
}

//go:nosplit
func i32x4_extmul_low_i16x8_u(a, b v128) (r v128) {
	for i := 0; i < 4; i++ {
		x := uint32(binary.LittleEndian.Uint16(a[2*i:])) * uint32(binary.LittleEndian.Uint16(b[2*i:]))
		binary.LittleEndian.PutUint32(r[4*i:], x)
	}
	return
}

//go:nosplit
func i32x4_extmul_high_i16x8_u(a, b v128) (r v128) {
	for i := 0; i < 4; i++ {
		x := uint32(binary.LittleEndian.Uint16(a[8+2*i:])) * uint32(binary.LittleEndian.Uint16(b[8+2*i:]))
		binary.LittleEndian.PutUint32(r[4*i:], x)
	}
	return
}

//go:nosplit
func i64x2_extmul_low_i32x4_s(a, b v128) (r v128) {
	for i := 0; i < 2; i++ {
		x := int64(int32(binary.LittleEndian.Uint32(a[4*i:]))) * int64(int32(binary.LittleEndian.Uint32(b[4*i:])))
		binary.LittleEndian.PutUint64(r[8*i:], uint64(x))
	}
	return
}

//go:nosplit
func i64x2_extmul_high_i32x4_s(a, b v128) (r v128) {
	for i := 0; i < 2; i++ {
		x := int64(int32(binary.LittleEndian.Uint32(a[8+4*i:]))) * int64(int32(binary.LittleEndian.Uint32(b[8+4*i:])))
		binary.LittleEndian.PutUint64(r[8*i:], uint64(x))
	}
	return
}

//go:nosplit
func i64x2_extmul_low_i32x4_u(a, b v128) (r v128) {
	for i := 0; i < 2; i++ {
		x := uint64(binary.LittleEndian.Uint32(a[4*i:])) * uint64(binary.LittleEndian.Uint32(b[4*i:]))
		binary.LittleEndian.PutUint64(r[8*i:], x)
	}
	return
}

//go:nosplit
func i64x2_extmul_high_i32x4_u(a, b v128) (r v128) {
	for i := 0; i < 2; i++ {
		x := uint64(binary.LittleEndian.Uint32(a[8+4*i:])) * uint64(binary.LittleEndian.Uint32(b[8+4*i:]))
		binary.LittleEndian.PutUint64(r[8*i:], x)
	}
	return
}

// Float unary operations.
// abs and neg operate on the sign bit, preserving NaN payloads.

//go:nosplit
func f32x4_abs(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], binary.LittleEndian.Uint32(a[i:])&0x7fffffff)
	}
	return
}

//go:nosplit
func f32x4_neg(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		binary.LittleEndian.PutUint32(r[i:], binary.LittleEndian.Uint32(a[i:])^0x80000000)
	}
	return
}

//go:nosplit
func f32x4_sqrt(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(float32(math.Sqrt(float64(x)))))
	}
	return
}

//go:nosplit
func f32x4_ceil(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(float32(math.Ceil(float64(x)))))
	}
	return
}

//go:nosplit
func f32x4_floor(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(float32(math.Floor(float64(x)))))
	}
	return
}

//go:nosplit
func f32x4_trunc(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(float32(math.Trunc(float64(x)))))
	}
	return
}

//go:nosplit
func f32x4_nearest(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(float32(math.RoundToEven(float64(x)))))
	}
	return
}

//go:nosplit
func f64x2_abs(a v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], binary.LittleEndian.Uint64(a[i:])&0x7fffffffffffffff)
	}
	return
}

//go:nosplit
func f64x2_neg(a v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		binary.LittleEndian.PutUint64(r[i:], binary.LittleEndian.Uint64(a[i:])^0x8000000000000000)
	}
	return
}

//go:nosplit
func f64x2_sqrt(a v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		binary.LittleEndian.PutUint64(r[i:], math.Float64bits(math.Sqrt(x)))
	}
	return
}

//go:nosplit
func f64x2_ceil(a v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		binary.LittleEndian.PutUint64(r[i:], math.Float64bits(math.Ceil(x)))
	}
	return
}

//go:nosplit
func f64x2_floor(a v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		binary.LittleEndian.PutUint64(r[i:], math.Float64bits(math.Floor(x)))
	}
	return
}

//go:nosplit
func f64x2_trunc(a v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		binary.LittleEndian.PutUint64(r[i:], math.Float64bits(math.Trunc(x)))
	}
	return
}

//go:nosplit
func f64x2_nearest(a v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		binary.LittleEndian.PutUint64(r[i:], math.Float64bits(math.RoundToEven(x)))
	}
	return
}

// Float arithmetic.

//go:nosplit
func f32x4_add(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(b[i:]))
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(x+y))
	}
	return
}

//go:nosplit
func f32x4_sub(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(b[i:]))
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(x-y))
	}
	return
}

//go:nosplit
func f32x4_mul(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(b[i:]))
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(x*y))
	}
	return
}

//go:nosplit
func f32x4_div(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(b[i:]))
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(x/y))
	}
	return
}

//go:nosplit
func f64x2_add(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		y := math.Float64frombits(binary.LittleEndian.Uint64(b[i:]))
		binary.LittleEndian.PutUint64(r[i:], math.Float64bits(x+y))
	}
	return
}

//go:nosplit
func f64x2_sub(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		y := math.Float64frombits(binary.LittleEndian.Uint64(b[i:]))
		binary.LittleEndian.PutUint64(r[i:], math.Float64bits(x-y))
	}
	return
}

//go:nosplit
func f64x2_mul(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		y := math.Float64frombits(binary.LittleEndian.Uint64(b[i:]))
		binary.LittleEndian.PutUint64(r[i:], math.Float64bits(x*y))
	}
	return
}

//go:nosplit
func f64x2_div(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		y := math.Float64frombits(binary.LittleEndian.Uint64(b[i:]))
		binary.LittleEndian.PutUint64(r[i:], math.Float64bits(x/y))
	}
	return
}

// Float min/max return a canonical NaN if either operand is NaN,
// and handle signed zeros: min(±0,∓0) is -0, max(±0,∓0) is +0.

//go:nosplit
func f32x4_min(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(b[i:]))
		var z uint32
		switch {
		case x != x || y != y:
			z = 0x7fc00000
		case x == y:
			z = math.Float32bits(x) | math.Float32bits(y)
		case x < y:
			z = math.Float32bits(x)
		default:
			z = math.Float32bits(y)
		}
		binary.LittleEndian.PutUint32(r[i:], z)
	}
	return
}

//go:nosplit
func f32x4_max(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(b[i:]))
		var z uint32
		switch {
		case x != x || y != y:
			z = 0x7fc00000
		case x == y:
			z = math.Float32bits(x) & math.Float32bits(y)
		case x > y:
			z = math.Float32bits(x)
		default:
			z = math.Float32bits(y)
		}
		binary.LittleEndian.PutUint32(r[i:], z)
	}
	return
}

//go:nosplit
func f64x2_min(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		y := math.Float64frombits(binary.LittleEndian.Uint64(b[i:]))
		var z uint64
		switch {
		case x != x || y != y:
			z = 0x7ff8000000000000
		case x == y:
			z = math.Float64bits(x) | math.Float64bits(y)
		case x < y:
			z = math.Float64bits(x)
		default:
			z = math.Float64bits(y)
		}
		binary.LittleEndian.PutUint64(r[i:], z)
	}
	return
}

//go:nosplit
func f64x2_max(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		y := math.Float64frombits(binary.LittleEndian.Uint64(b[i:]))
		var z uint64
		switch {
		case x != x || y != y:
			z = 0x7ff8000000000000
		case x == y:
			z = math.Float64bits(x) & math.Float64bits(y)
		case x > y:
			z = math.Float64bits(x)
		default:
			z = math.Float64bits(y)
		}
		binary.LittleEndian.PutUint64(r[i:], z)
	}
	return
}

// Pseudo-minimum/maximum: pmin is y < x ? y : x, pmax is x < y ? y : x.

//go:nosplit
func f32x4_pmin(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(b[i:]))
		if y < x {
			binary.LittleEndian.PutUint32(r[i:], math.Float32bits(y))
		} else {
			binary.LittleEndian.PutUint32(r[i:], math.Float32bits(x))
		}
	}
	return
}

//go:nosplit
func f32x4_pmax(a, b v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		y := math.Float32frombits(binary.LittleEndian.Uint32(b[i:]))
		if x < y {
			binary.LittleEndian.PutUint32(r[i:], math.Float32bits(y))
		} else {
			binary.LittleEndian.PutUint32(r[i:], math.Float32bits(x))
		}
	}
	return
}

//go:nosplit
func f64x2_pmin(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		y := math.Float64frombits(binary.LittleEndian.Uint64(b[i:]))
		if y < x {
			binary.LittleEndian.PutUint64(r[i:], math.Float64bits(y))
		} else {
			binary.LittleEndian.PutUint64(r[i:], math.Float64bits(x))
		}
	}
	return
}

//go:nosplit
func f64x2_pmax(a, b v128) (r v128) {
	for i := 0; i < 16; i += 8 {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[i:]))
		y := math.Float64frombits(binary.LittleEndian.Uint64(b[i:]))
		if x < y {
			binary.LittleEndian.PutUint64(r[i:], math.Float64bits(y))
		} else {
			binary.LittleEndian.PutUint64(r[i:], math.Float64bits(x))
		}
	}
	return
}

// Float/integer conversions.

//go:nosplit
func f32x4_demote_f64x2_zero(a v128) (r v128) {
	for i := 0; i < 2; i++ {
		x := math.Float64frombits(binary.LittleEndian.Uint64(a[8*i:]))
		binary.LittleEndian.PutUint32(r[4*i:], math.Float32bits(float32(x)))
	}
	return
}

//go:nosplit
func f64x2_promote_low_f32x4(a v128) (r v128) {
	for i := 0; i < 2; i++ {
		x := math.Float32frombits(binary.LittleEndian.Uint32(a[4*i:]))
		binary.LittleEndian.PutUint64(r[8*i:], math.Float64bits(float64(x)))
	}
	return
}

//go:nosplit
func f32x4_convert_i32x4_s(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := int32(binary.LittleEndian.Uint32(a[i:]))
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(float32(x)))
	}
	return
}

//go:nosplit
func f32x4_convert_i32x4_u(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		x := binary.LittleEndian.Uint32(a[i:])
		binary.LittleEndian.PutUint32(r[i:], math.Float32bits(float32(x)))
	}
	return
}

//go:nosplit
func f64x2_convert_low_i32x4_s(a v128) (r v128) {
	for i := 0; i < 2; i++ {
		x := int32(binary.LittleEndian.Uint32(a[4*i:]))
		binary.LittleEndian.PutUint64(r[8*i:], math.Float64bits(float64(x)))
	}
	return
}

//go:nosplit
func f64x2_convert_low_i32x4_u(a v128) (r v128) {
	for i := 0; i < 2; i++ {
		x := binary.LittleEndian.Uint32(a[4*i:])
		binary.LittleEndian.PutUint64(r[8*i:], math.Float64bits(float64(x)))
	}
	return
}

//go:nosplit
func i32x4_trunc_sat_f32x4_s(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		f := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		var x int32
		switch {
		case f <= math.MinInt32:
			x = math.MinInt32
		case f >= math.MaxInt32:
			x = math.MaxInt32
		case f != f:
			x = 0
		default:
			x = int32(f)
		}
		binary.LittleEndian.PutUint32(r[i:], uint32(x))
	}
	return
}

//go:nosplit
func i32x4_trunc_sat_f32x4_u(a v128) (r v128) {
	for i := 0; i < 16; i += 4 {
		f := math.Float32frombits(binary.LittleEndian.Uint32(a[i:]))
		var x uint32
		switch {
		case f <= 0 || f != f:
			x = 0
		case f >= math.MaxUint32:
			x = math.MaxUint32
		default:
			x = uint32(f)
		}
		binary.LittleEndian.PutUint32(r[i:], x)
	}
	return
}

//go:nosplit
func i32x4_trunc_sat_f64x2_s_zero(a v128) (r v128) {
	for i := 0; i < 2; i++ {
		f := math.Float64frombits(binary.LittleEndian.Uint64(a[8*i:]))
		var x int32
		switch {
		case f <= math.MinInt32:
			x = math.MinInt32
		case f >= math.MaxInt32:
			x = math.MaxInt32
		case f != f:
			x = 0
		default:
			x = int32(f)
		}
		binary.LittleEndian.PutUint32(r[4*i:], uint32(x))
	}
	return
}

//go:nosplit
func i32x4_trunc_sat_f64x2_u_zero(a v128) (r v128) {
	for i := 0; i < 2; i++ {
		f := math.Float64frombits(binary.LittleEndian.Uint64(a[8*i:]))
		var x uint32
		switch {
		case f <= 0 || f != f:
			x = 0
		case f >= math.MaxUint32:
			x = math.MaxUint32
		default:
			x = uint32(f)
		}
		binary.LittleEndian.PutUint32(r[4*i:], x)
	}
	return
}
