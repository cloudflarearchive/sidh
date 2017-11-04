
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

----
## Benchmark Comparison
For Shared Secret computation, the execution time was reduced by 9.5%.


| benchmark                         |   old ns/op  |   new ns/op   |  delta |
|---------------------------------------------------------------------------|
| BenchmarkAliceKeyGen-4            |   33945573   |   32551477    | -4.11% |
| BenchmarkAliceKeyGenSlow-4        |   289292459  |   286781846   | -0.87% |
| BenchmarkBobKeyGen-4              |   37633036   |   37420870    | -0.56% |
| BenchmarkBobKeyGenSlow-4          |   472687918  |   465635415   | -1.49% |
| BenchmarkSharedSecretAlice-4      |   32414528   |   29297806    | -9.62% |
| BenchmarkSharedSecretAliceSlow-4  |   287817189  |   287222653   | -0.21% |
| BenchmarkSharedSecretBob-4        |   37216521   |   33594430    | -9.73% |
| BenchmarkSharedSecretBobSlow-4    |   473442424  |   464861983   | -1.81% |

----
