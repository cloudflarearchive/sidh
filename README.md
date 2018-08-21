# `p751sidh`

Project provides a Go implementation of  (ephemeral) supersingular isogeny Diffie-Hellman (SIDH) and supersingular isogeny key exchange (SIKE), as specified in [SIDH-spec, PQC NIST Submission](http://sike.org/files/SIDH-spec.pdf) (Nov 30, 2017).

The implementation is intended for use on the `amd64` architecture only -- no
generic field arithmetic implementation is provided.  Portions of the field
arithmetic were ported from the Microsoft Research implementation.

The SIDH package does NOT implement key validation. It means that it should only be
used for ephemeral DH. Each keypair should be used at most once.

If you feel that SIDH may be appropriate for you, consult your cryptographer.

## Source code
Project provides following packages:
* ``p751toolbox``: P751 field arithmetic, curve computation and isogeny internal functions
* ``sidh``: Implementation of SIDH key agreement
* ``sike``: Implementation of SIKE PKE and KEM, based on ``sidh`` package

## Testing
At development time following make targets may come handy:
* ``make test`` : unit testing
* ``make bench``: benchmarking
* ``make cover``: produces code coverage as txt file (used by travis)

It is possible to add one of following postfixes to each of the targets above, in order to run tests specific to ``-p751toolbox``, ``-sidh`` or ``-sike``.

## Acknowledgements

Special thanks to [Craig Costello](http://www.craigcostello.com.au/), [Diego Aranha](https://sites.google.com/site/dfaranha/), and [Deirdre Connolly](https://twitter.com/durumcrustulum) for advice
and discussion.

