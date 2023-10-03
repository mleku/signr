# signr

Designed along similar lines to `ssh-keygen` but with a more singular purpose,
`signr` can generate new keys, keeping them in a directory with a familiar
similar layout as `.ssh`, as well as sign and verify.

It provides the following functionality:

- [ ] Key generation - using the system's strong entropy source

- Key import

    - [ ] hexadecimal
    - [ ] nsec
    - [ ] Bech32
    - [ ] BIP39 word-keys

- [ ] Signing - using a distinct protocol to keep the signature space
  isolated from other protocols, such as Bitcoin message signatures. The 
  signature algorithm, however, will be the standard ECDSA as used in 
  Bitcoin as this compactly includes 

- [ ] Verification - checking that a signature matches a given file or hash on
  a file

In order to prevent cross-protocol attacks, the signature is applied not
directly on the hash of the message, but rather a distinctive structure
that prevents collisions between a signature for this versus any other use
of the same hash functions and elliptic curves.

The raw bytes that are hashed using SHA256 are constructed as follows:

1. Magic - the string `signr`. Sometimes referred to as a "namespace".
2. Version - Monotonic number as a string encoding the version being used.
   Starts with 0.
3. The hash function used, `SHA256` normally but allowing future additions,
   This is also part of the signature prefix.
4. The signature algorithm is here encoded. It can be ECDSA for standard 
   Bitcoin transaction signatures, which validate by producing the public 
   key, thus eliminating the need to additionally specify the public key to 
   the validator, allowing this key to be searched instead of pre-specified.
   SCHNORR can be used to enable standard btcec Schnorr signatures, these 
   require the verifier to also know the public key, and yield only a 
   boolean result. ECDSA recovered signature yields the boolean by comparing 
   with the provided key.
5. Nonce - a strong random value of 64 bits as 16 hex
   characters, that are repeated in the signature prefix to enable the
   generation of the actual message hash that is signed
6. Hash of the message being signed in hex

The string is interpreted as standard ASCII, and the hash that is signed is
generated from these ASCII bytes. All hexadecimal digits are lower case. 
Hash function strings should be as they are normally written as identifiers, 
usually all caps, this is necessary otherwise there could be malleability.
Each section is separated by a underscore, so the whole string is selected 
"by word" with most GUI text selection systems using a double click.

The canonical encoding of the signature prefix would thus look like this:

   signr_0_SHA256_ECDSA_deadbeefcafeb00b_0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef

The signature will then be in Bech32, and appended in place of the 
hexadecimal hash string as shown above, and the verifier must first generate 
the message hash, hash this string, and then validate it against the decoded 
Bech32 signature that was in the last place in the published signature.
