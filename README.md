# signr

Designed along similar lines to `ssh-keygen`, `signr` can generate new keys, password protect them, keeping them in a directory with a familiar
similar layout as `.ssh`, as well as sign and verify.

It provides the following functionality:

- [x] Key generation - using the system's strong entropy source

    - [x] hexadecimal
    - [x] nsec

- Secret key import 

    - [x] hexadecimal
    - [x] nsec

- [ ] Signing - using a distinct protocol to keep the signature space
  isolated from other protocols, such as Bitcoin message signatures.

- [ ] Verification - checking that a signature matches a given file or hash on
  a file

- [ ] Keychain management - storing keys in user profile and validating security of these files (configuration and access similar to ssh with common tools).

- [ ] Encryption of private keys.

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
   generation of the actual message hash that is signed. This is here to 
   ensure the same hash is never signed twice as this weakens the security 
   of the EC keys.
6. Public Key of the signatory, required for verification
7. Hash of the message being signed in hex

The string is interpreted as standard ASCII, and the hash that is signed is
generated from these ASCII bytes. All hexadecimal digits are lower case. 
Hash function strings should be as they are normally written as identifiers, 
usually all caps, this is necessary otherwise there could be malleability.
Each section is separated by a underscore, so the whole string is selected 
"by word" with most GUI text selection systems using a double click.

The canonical encoding of the signature prefix would thus look like this:

    signr_0_SHA256_ECDSA_deadbeefcafeb00b_npub1e44x0gq7xg2rln2ffyy4ck5ghyt03mstacupksjy462u50nqux6qt8zpf8_0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef

The provided signature contains the Bech32 encoded signature bytes in the 
place of the message hash. The verifier splits off this signature, adds the 
message hash in its place, hashes the resulting string, and then after decoding
the signature to bytes, calls the secp256k1 schnorr signature verify function 
and gets the result. If the signature was made on the hash of this same string 
then it will pass.