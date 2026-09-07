;; Regression: a provided import may call a helper (here load64) that the
;; module's own code never uses; the translator must still emit it.
(module
  (import "env" "peek" (func $peek (param i32) (result i64)))
  (memory (export "memory") 1)
  (data (i32.const 8) "\01\02\03\04\05\06\07\08")
  (func (export "test") (result i64)
    i32.const 8
    call $peek))
