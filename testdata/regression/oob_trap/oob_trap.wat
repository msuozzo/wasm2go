;; Regression: linear-memory accesses must trap deterministically at the
;; memory boundary, for every access width, with and without static offsets,
;; before and after memory.grow. Guards the single-bounds-check emission used
;; by -unsafe for 32-bit memories.
(module
  (memory (export "memory") 1)
  (func (export "ld16") (param i32) (result i32)
    local.get 0 i32.load16_u)
  (func (export "ld32") (param i32) (result i32)
    local.get 0 i32.load)
  (func (export "ld64") (param i32) (result i64)
    local.get 0 i64.load)
  (func (export "ld32o") (param i32) (result i32)
    local.get 0 i32.load offset=0xffffffff)
  (func (export "st32") (param i32 i32)
    local.get 0 local.get 1 i32.store)
  (func (export "st64") (param i32 i64)
    local.get 0 local.get 1 i64.store)
  (func (export "grow") (param i32) (result i32)
    local.get 0 memory.grow)
)
