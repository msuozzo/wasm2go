// Native SIMD helper overrides using the Go 1.27 experimental simd/archsimd
// package. This file is parsed by wasm2go, never compiled as part of this
// package: under -simd-only these bodies are emitted instead of the portable
// ones, for goexperiment.simd && (amd64 || arm64).
//
// Only ops whose archsimd lowering is Wasm-faithful on BOTH amd64 (within
// AVX/AVX2) and arm64 (Neon) are included; everything else falls back to the
// portable helpers. Deliberately portable: float min/max (Wasm NaN and signed
// zero rules), pmin/pmax, i8x16 shifts and i64x2.shr_s (no x86 encoding below
// AVX512), i64x2 mul/abs, popcnt, narrows and high extends, saturating float
// to int truncations (x86 gives INT_MIN, Wasm saturates), extadd/extmul/dot/
// q15mulr, shuffles, lane and memory access, and mask reductions.

//go:build ignore

package helpers

import "simd/archsimd"

//go:nosplit
func i8x16_splat(x int32) (r v128) {
	archsimd.BroadcastUint8x16(uint8(x)).Store(r[:])
	return
}

//go:nosplit
func i16x8_splat(x int32) (r v128) {
	archsimd.BroadcastUint16x8(uint16(x)).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_splat(x int32) (r v128) {
	archsimd.BroadcastUint32x4(uint32(x)).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_splat(x int64) (r v128) {
	archsimd.BroadcastUint64x2(uint64(x)).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_splat(x float32) (r v128) {
	archsimd.BroadcastFloat32x4(x).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_splat(x float64) (r v128) {
	archsimd.BroadcastFloat64x2(x).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func v128_not(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	x.Not().Store(r[:])
	return
}

//go:nosplit
func v128_and(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.And(y).Store(r[:])
	return
}

//go:nosplit
func v128_andnot(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.AndNot(y).Store(r[:])
	return
}

//go:nosplit
func v128_or(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.Or(y).Store(r[:])
	return
}

//go:nosplit
func v128_xor(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.Xor(y).Store(r[:])
	return
}

//go:nosplit
func i8x16_add(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.Add(y).Store(r[:])
	return
}

//go:nosplit
func i8x16_sub(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.Sub(y).Store(r[:])
	return
}

//go:nosplit
func i16x8_add(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.Add(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_sub(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.Sub(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_add(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.Add(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_sub(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.Sub(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_add(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s()
	x.Add(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_sub(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s()
	x.Sub(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_mul(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.Mul(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_mul(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.Mul(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i8x16_add_sat_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	y := archsimd.LoadUint8x16(b[:]).BitsToInt8()
	x.AddSaturated(y).ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_add_sat_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.AddSaturated(y).Store(r[:])
	return
}

//go:nosplit
func i8x16_sub_sat_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	y := archsimd.LoadUint8x16(b[:]).BitsToInt8()
	x.SubSaturated(y).ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_sub_sat_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.SubSaturated(y).Store(r[:])
	return
}

//go:nosplit
func i8x16_avgr_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.Average(y).Store(r[:])
	return
}

//go:nosplit
func i16x8_add_sat_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s().BitsToInt16()
	x.AddSaturated(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_add_sat_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.AddSaturated(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_sub_sat_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s().BitsToInt16()
	x.SubSaturated(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_sub_sat_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.SubSaturated(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_avgr_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.Average(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i8x16_min_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	y := archsimd.LoadUint8x16(b[:]).BitsToInt8()
	x.Min(y).ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_min_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.Min(y).Store(r[:])
	return
}

//go:nosplit
func i8x16_max_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	y := archsimd.LoadUint8x16(b[:]).BitsToInt8()
	x.Max(y).ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_max_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.Max(y).Store(r[:])
	return
}

//go:nosplit
func i8x16_abs(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	x.Abs().ToBits().Store(r[:])
	return
}

//go:nosplit
func i16x8_min_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s().BitsToInt16()
	x.Min(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_min_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.Min(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_max_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s().BitsToInt16()
	x.Max(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_max_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.Max(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_abs(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	x.Abs().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_min_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToInt32()
	x.Min(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_min_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.Min(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_max_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToInt32()
	x.Max(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_max_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.Max(y).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_abs(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	x.Abs().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i8x16_neg(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	x.Neg().ToBits().Store(r[:])
	return
}

//go:nosplit
func i16x8_neg(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	x.Neg().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_neg(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	x.Neg().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_neg(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToInt64()
	x.Neg().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_shl(a v128, y int32) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	x.ShiftAllLeft(uint64(y) & 15).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_shr_u(a v128, y int32) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	x.ShiftAllRight(uint64(y) & 15).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_shr_s(a v128, y int32) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	x.ShiftAllRight(uint64(y) & 15).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_shl(a v128, y int32) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	x.ShiftAllLeft(uint64(y) & 31).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_shr_u(a v128, y int32) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	x.ShiftAllRight(uint64(y) & 31).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_shr_s(a v128, y int32) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	x.ShiftAllRight(uint64(y) & 31).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_shl(a v128, y int32) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s()
	x.ShiftAllLeft(uint64(y) & 63).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_shr_u(a v128, y int32) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s()
	x.ShiftAllRight(uint64(y) & 63).ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i8x16_eq(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.Equal(y).ToInt8x16().ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_ne(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.NotEqual(y).ToInt8x16().ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_lt_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	y := archsimd.LoadUint8x16(b[:]).BitsToInt8()
	x.Less(y).ToInt8x16().ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_lt_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.Less(y).ToInt8x16().ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_gt_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	y := archsimd.LoadUint8x16(b[:]).BitsToInt8()
	x.Greater(y).ToInt8x16().ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_gt_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.Greater(y).ToInt8x16().ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_le_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	y := archsimd.LoadUint8x16(b[:]).BitsToInt8()
	x.LessEqual(y).ToInt8x16().ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_le_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.LessEqual(y).ToInt8x16().ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_ge_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	y := archsimd.LoadUint8x16(b[:]).BitsToInt8()
	x.GreaterEqual(y).ToInt8x16().ToBits().Store(r[:])
	return
}

//go:nosplit
func i8x16_ge_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	y := archsimd.LoadUint8x16(b[:])
	x.GreaterEqual(y).ToInt8x16().ToBits().Store(r[:])
	return
}

//go:nosplit
func i16x8_eq(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.Equal(y).ToInt16x8().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_ne(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.NotEqual(y).ToInt16x8().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_lt_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s().BitsToInt16()
	x.Less(y).ToInt16x8().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_lt_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.Less(y).ToInt16x8().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_gt_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s().BitsToInt16()
	x.Greater(y).ToInt16x8().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_gt_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.Greater(y).ToInt16x8().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_le_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s().BitsToInt16()
	x.LessEqual(y).ToInt16x8().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_le_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.LessEqual(y).ToInt16x8().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_ge_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s().BitsToInt16()
	x.GreaterEqual(y).ToInt16x8().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_ge_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint16s()
	x.GreaterEqual(y).ToInt16x8().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_eq(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.Equal(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_ne(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.NotEqual(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_lt_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToInt32()
	x.Less(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_lt_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.Less(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_gt_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToInt32()
	x.Greater(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_gt_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.Greater(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_le_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToInt32()
	x.LessEqual(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_le_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.LessEqual(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_ge_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToInt32()
	x.GreaterEqual(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_ge_u(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s()
	x.GreaterEqual(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_eq(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToInt64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToInt64()
	x.Equal(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_ne(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToInt64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToInt64()
	x.NotEqual(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_lt_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToInt64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToInt64()
	x.Less(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_gt_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToInt64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToInt64()
	x.Greater(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_le_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToInt64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToInt64()
	x.LessEqual(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_ge_s(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToInt64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToInt64()
	x.GreaterEqual(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_add(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToFloat32()
	x.Add(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_sub(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToFloat32()
	x.Sub(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_mul(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToFloat32()
	x.Mul(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_div(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToFloat32()
	x.Div(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_sqrt(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	x.Sqrt().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_abs(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	x.Abs().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_neg(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	x.Neg().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_ceil(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	x.Ceil().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_floor(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	x.Floor().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_trunc(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	x.Trunc().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_nearest(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	x.Round().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_eq(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToFloat32()
	x.Equal(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_ne(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToFloat32()
	x.NotEqual(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_lt(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToFloat32()
	x.Less(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_gt(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToFloat32()
	x.Greater(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_le(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToFloat32()
	x.LessEqual(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_ge(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToFloat32()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint32s().BitsToFloat32()
	x.GreaterEqual(y).ToInt32x4().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_add(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToFloat64()
	x.Add(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_sub(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToFloat64()
	x.Sub(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_mul(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToFloat64()
	x.Mul(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_div(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToFloat64()
	x.Div(y).ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_sqrt(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	x.Sqrt().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_abs(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	x.Abs().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_neg(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	x.Neg().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_ceil(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	x.Ceil().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_floor(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	x.Floor().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_trunc(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	x.Trunc().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_nearest(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	x.Round().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_eq(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToFloat64()
	x.Equal(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_ne(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToFloat64()
	x.NotEqual(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_lt(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToFloat64()
	x.Less(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_gt(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToFloat64()
	x.Greater(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_le(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToFloat64()
	x.LessEqual(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f64x2_ge(a, b v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	y := archsimd.LoadUint8x16(b[:]).ReshapeToUint64s().BitsToFloat64()
	x.GreaterEqual(y).ToInt64x2().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_extend_low_i8x16_s(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).BitsToInt8()
	x.ExtendLo8ToInt16().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i16x8_extend_low_i8x16_u(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:])
	x.ExtendLo8ToUint16().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_extend_low_i16x8_s(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s().BitsToInt16()
	x.ExtendLo4ToInt32().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i32x4_extend_low_i16x8_u(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint16s()
	x.ExtendLo4ToUint32().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_extend_low_i32x4_s(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	x.ExtendLo2ToInt64().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func i64x2_extend_low_i32x4_u(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s()
	x.ExtendLo2ToUint64().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_demote_f64x2_zero(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint64s().BitsToFloat64()
	x.ConvertToFloat32().ToBits().ReshapeToUint8s().Store(r[:])
	return
}

//go:nosplit
func f32x4_convert_i32x4_s(a v128) (r v128) {
	x := archsimd.LoadUint8x16(a[:]).ReshapeToUint32s().BitsToInt32()
	x.ConvertToFloat32().ToBits().ReshapeToUint8s().Store(r[:])
	return
}
