
## List of Optimizations

---
### General optimizations

The following are general and well-known optimizations.

1. Replace ADD+DOUBLE by a xDBLADD function.
  This reduces 1M per iteration in 3-point ladder and 2M per iteration in ScalarMult function.
2. Print for types (useful for debugging purposes)


----
### Optimizations derived from FLOR-SIDH-x64
The following are specific optimizations based on the FLOR-SIDH-x64 work.

1. New formulas for point tripling.
2. New R2L method for shared secret, replaces 3-point ladder.

