# signr

Designed along similar lines to `ssh-keygen`, `signr` can generate new keys, password protect them, keeping them in a directory with a familiar
similar layout as `.ssh`, as well as sign and verify.

It provides the following functionality:

- [x] Key generation - using the system's strong entropy source

    - [x] hexadecimal
    - [x] nsec
- [x] Secret key import 
    - [x] hexadecimal
    - [x] nostr nsec key format
- [ ] Signing of files - via path or via stdin
- [ ] Verification - checking a signature matches a provided file or stream from stdin
  
- [ ] Keychain management 
    - [x]  storing keys in user profile 

    - [ ] validating filesystem security of these files 

    - [x] setting a default key to use when unspecified for signing

    - [x] Encryption of private keys.


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
4. The signature algorithm is here encoded. It supports BIP-340 style 32 byte public keys as specified in NIP-19, and produces 64 byte Schnorr signatures. It is specified to enable future expansion to support other signature algorithms.
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