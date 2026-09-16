#pragma once

#define __STDC_ENDIAN_BIG__ 4321
#define __STDC_ENDIAN_LITTLE__ 1234
#define __STDC_ENDIAN_NATIVE__ __STDC_ENDIAN_LITTLE__

#define stdc_bit_ceil(x) (__builtin_stdc_bit_ceil(x))
#define stdc_bit_floor(x) (__builtin_stdc_bit_floor(x))
#define stdc_bit_width(x) (__builtin_stdc_bit_width(x))
#define stdc_count_ones(x) (__builtin_stdc_count_ones(x))
#define stdc_count_zeros(x) (__builtin_stdc_count_zeros(x))
#define stdc_first_leading_one(x) (__builtin_stdc_first_leading_one(x))
#define stdc_first_leading_zero(x) (__builtin_stdc_first_leading_zero(x))
#define stdc_first_trailing_one(x) (__builtin_stdc_first_trailing_one(x))
#define stdc_first_trailing_zero(x) (__builtin_stdc_first_trailing_zero(x))
#define stdc_has_single_bit(x) (__builtin_stdc_has_single_bit(x))
#define stdc_leading_ones(x) (__builtin_stdc_leading_ones(x))
#define stdc_leading_zeros(x) (__builtin_stdc_leading_zeros(x))
#define stdc_rotate_left(x, y) (__builtin_stdc_rotate_left(x, y))
#define stdc_rotate_right(x, y) (__builtin_stdc_rotate_right(x, y))
#define stdc_trailing_ones(x) (__builtin_stdc_trailing_ones(x))
#define stdc_trailing_zeros(x) (__builtin_stdc_trailing_zeros(x))
