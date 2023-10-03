# signostr

Designed along similar lines to `ssh-keygen` but with a more singular purpose,
`signostr` can generate new keys, keeping them in a directory with a familiar
similar layout as `.ssh`, as well as sign and verify.

It provides the following functionality:

- [ ] Key generation - using the system's strong entropy source, or

- Key import

    - [ ] hexadecimal
    - [ ] nsec
    - [ ] Bech32
    - [ ] BIP39 word-keys

- [ ] Signing - using a distinct protocol to keep the signature space
  isolated from other protocols, such as Bitcoin message signatures

- [ ] Verification - checking that a signature matches a given file or hash on
  a file

In order to prevent cross-protocol attacks, the signature is applied not
directly on the hash of the message, but rather a distinctive structure
that prevents collisions between a signature for this versus any other use
of the same hash functions and elliptic curves.

The raw bytes that are hashed using SHA256 are constructed as follows:

1. Magic - the string `signostr`. Sometimes referred to as a "namespace".
2. Version - Monotonic number as a string encoding the version being used.
   Starts with 0.
3. The hash function used, `SHA256` normally but allowing future additions,
   This is also part of the signature prefix.
4. Nonce - a strong random value of 64 bits as 16 hex
   characters, that are repeated in the signature prefix to enable the
   generation of the actual message hash that is signed
5. Hash of the message being signed in hex

The string is interpreted as standard ASCII, and the hash that is signed is
generated from these ASCII bytes.