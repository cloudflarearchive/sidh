# `sidh`

Project provides a Go implementation of  (ephemeral) supersingular isogeny Diffie-Hellman (SIDH) and supersingular isogeny key exchange (SIKE), as specified in [SIDH-spec, PQC NIST Submission](http://sike.org/files/SIDH-spec.pdf) (Nov 30, 2017).

Portions of the field arithmetic were ported from the Microsoft Research implementation.

The SIDH package does NOT implement key validation. It means that it should only be
used for ephemeral DH. Each keypair should be used at most once.

If you feel that SIDH may be appropriate for you, consult your cryptographer.

## Source code
Project provides following packages:
* ``p503``: P503 field arithmetic
* ``p751``: P751 field arithmetic
* ``sidh``: Implementation of SIDH key agreement
* ``sike``: Implementation of SIKE PKE and KEM, based on ``sidh`` package

## Testing
At development time following make targets may come handy:
* ``make test`` : unit testing
* ``make bench``: benchmarking
* ``make cover``: produces code coverage as txt file (used by travis)

It is possible to add one of following postfixes to each of the targets above, in order to run tests specific to ``-p503``, ``-p751``, ``-sidh`` or ``-sike``.

## Acknowledgements

Special thanks to [Craig Costello](http://www.craigcostello.com.au/), [Diego Aranha](https://sites.google.com/site/dfaranha/), and [Deirdre Connolly](https://twitter.com/durumcrustulum) for advice
and discussion.

