# signr

Designed along similar lines to `ssh-keygen`, `signr` can generate new keys,
password protect them, keeping them in a directory with a familiar similar
layout as `.ssh`, as well as sign and verify.

## usage

### help

The first thing to know is how to access the inbuilt help. The `signr` command by default prints the top level help information with no arguments:

```bash
# signr
```

will print something like this:

    A command line interface for generating, importing, signing, verifying and managing keys used with the Nostr protocol.
    
    Designed to function in a similar way to ssh-keygen in that it keeps the keychain in a user directory with named key pairs and a configuration file.
    
    Usage:
      signr [command]
    
    Available Commands:
      completion  Generate the autocompletion script for the specified shell
      gen         Generate a new nostr key
      help        Help about any command
      import      Import a secret key
      listkeys    List the keys in the keychain
      set         Set configuration values from the CLI
      sign        Generate a signature on a file
      verify      check that a file matches a signature
    
    Flags:
      -h, --help      help for signr
      -v, --verbose   prints more things
    
    Use "signr [command] --help" for more information about a command.

Where it makes sense, the `-v` or `--verbose` flag prints additional information, usually to `stderr`, where it won't be caught by a `stdout` pipe in scripts, but will print information for the user to read.

When requesting help information, you can either use the top level command `help` as in this:

    signr help set

Or you can use the flag version:

    signr set --help

This applies no matter how deep you go, such as the `set default` command:

    signr help set default

or

    signr set default --help

The subcommand version (first) is more intuitive and fluid especially if you are a touch typist.

### key generation

    signr gen newkeyname

A name is required when generating a new key. As contrasted with the interface of ssh, where naming new keys is optional, it is mandatory in `signr` because we want users to think of this as a keychain in the first class sense, and to be able to use it this way.

The first key generated will also become the default key, which can be changed using the `set` command.

After running gen, you will now have three files in a new subdirectory `~/.signr`:

    config.yaml
    newkeyname
    newkeyname.pub

`config.yaml` will have the following contents:

    default: newkeyname

If you use the same name in this command as an existing key, it will refuse to create a key, explaining that the name is taken.

A delete command is provided, that renames the key with a random additional string before an added extension `.del` and if you REALLY must delete it, you can use `rm ~/.signr/keyname.123abc.del.*` (todo: wip)

### set default

The first key is default, and this enables you to omit the key name when performing sign operations. Once you add more keys, you may want to change the default key, and to do that:

    signr set default newdefaultkey

This checks the `newdefaultkey` exists and then modifies the `config.yaml` to reflect the new status. After this, sign operations without a key named will use this key.

### sign

To sign a file or hash using a key from the keychain. use the following command

    signr sign <filename|-|hash> [key name]

If the first parameter is `-` signr assumes the file to be hashed will appear on the standard input, IE: via a pipe, eg:

    cat filename | signr sign - [key name]

Either of these options will produce a canonical binary blob/file signature thus:

    $ go run . sign go.sum
    
    signr_0_SHA256_SCHNORR_088b9a5b4bf3fd87_npub1vp0q3vl6sgvpwq5pcfdjyjj8q2kg3r6aa45vclzu8xuyrmaahl4svlvppa_sig135x8pmhfe2ancwjxeeu3rfdumtr95ceaqcw2vx2nv2ek655hyq5unc4g968r8jt4ucflqsz7m2l9jh2dy36va5kr4nzrgpgvrt2gzjcph9zgu

Note the presence of the random nonce and the public key.

For other uses, it may be desired to only sign a hash and only get the signature back.

If the first parameter is 64 characters long, and purely containing `0-9a-f` it will be interpreted as a hash and signed without the addition of the random number and public key as with the `filename/-` and will return only the signature in bech32.

    $ signr sign 5167fc037d96be70b809e3c669b7cb64864206273757f670af8f729beabf6ddc
    
    sig1qvltpcp06r902gwmylvysfdqh453c2fn9a9znr0y5jvcsm7j75snjcqxyhvr8n4hhwdcw6ex0w8g24nl8km5mrmjg0sw0r889unfvsslqe2rv

The signature will be printed in Bech32 by default, use the flag `--hex` to get a hexadecimal version:

    $ signr sign --hex 5167fc037d96be70b809e3c669b7cb64864206273757f670af8f729beabf6ddc
    
    033eb0e02fd0caf521db27d84825a0bd691c29332f4a298de4a499886fd2f52139600625d833ceb7bb9b876b267b8e85567f3db74d8f7243e0e78ce72f269642


For reasons of prevention of namespace attacks, these signatures are generated on a structured value that is hashed, like so:

    signr_0_SHA256_SCHNORR_5167fc037d96be70b809e3c669b7cb64864206273757f670af8f729beabf6ddc

rather than being signed directly on the provided value encoded as binary.

To allow custom protocols to be devised, this can have an additional custom string using the `--custom` flag:

    $ go run . -v sign --custom arbitrary-protocol-string go.sum
    Using config file: /home/me/.signr/config.yaml
    signing on message: signr_0_SHA256_SCHNORR_arbitrary-protocol-string_8b716747195afc47_npub1vp0q3vl6sgvpwq5pcfdjyjj8q2kg3r6aa45vclzu8xuyrmaahl4svlvppa_a61e297a6a401ceef85dff780e757d55016e402230a8c5e308fcb13cf4434f3a
    signr_0_SHA256_SCHNORR_arbitrary-protocol-string_8b716747195afc47_npub1vp0q3vl6sgvpwq5pcfdjyjj8q2kg3r6aa45vclzu8xuyrmaahl4svlvppa_sig1wk3ndpaldxct9v9dc3p0r0ynpwzcy9nyuadyzfp7xx7hcrxrm6c4ptanwlqywnf0vdlyc28krnfszyw0apev2t5vyrwetgp6deqvc7ggkf4xy

The arbitrary custom string has any leading or following spaces or carriage returns removed:

    $ go run . -v sign --custom "arbitrary protocol string 
    " go.sum
    Using config file: /home/me/.signr/config.yaml
    signing on message: signr_0_SHA256_SCHNORR_arbitrary_protocol_string_1093f338208d86c6_npub1vp0q3vl6sgvpwq5pcfdjyjj8q2kg3r6aa45vclzu8xuyrmaahl4svlvppa_a61e297a6a401ceef85dff780e757d55016e402230a8c5e308fcb13cf4434f3a
    signr_0_SHA256_SCHNORR_arbitrary_protocol_string_1093f338208d86c6_npub1vp0q3vl6sgvpwq5pcfdjyjj8q2kg3r6aa45vclzu8xuyrmaahl4svlvppa_sig1f8w9pgwuuj9xhg023g46xjtuz9h958rngy9mt2zqku9r0aj9cfg8rst4dxccqzdvvum9m3hu3ulhxfy6wsspgpvujgr4gt9dmatkyusn2ldld

Carriage returns as well, though these can only end up in the text very deliberately with a copy/paste or other programmatic method:

    $ go run . -v sign --custom "arbitrary protocol string 
    " go.sum
    Using config file: /home/me/.signr/config.yaml
    signing on message: signr_0_SHA256_SCHNORR_arbitrary_protocol_string_93d2fdebb8c517da_npub1vp0q3vl6sgvpwq5pcfdjyjj8q2kg3r6aa45vclzu8xuyrmaahl4svlvppa_a61e297a6a401ceef85dff780e757d55016e402230a8c5e308fcb13cf4434f3a
    signr_0_SHA256_SCHNORR_arbitrary_protocol_string_93d2fdebb8c517da_npub1vp0q3vl6sgvpwq5pcfdjyjj8q2kg3r6aa45vclzu8xuyrmaahl4svlvppa_sig189c3ty6nqjg2v8zalgej26vq2sms7n8hku5kvvc8v2whveas4n8vzpmfk5xx64ma7h80eq0czj2f69yrfavyet9vljwkns40r0eqpqgwg7urt

For completeness, the same thing except with a hash instead of a file:

    $ go run . -v sign --custom "arbitrary protocol string 
    " a61e297a6a401ceef85dff780e757d55016e402230a8c5e308fcb13cf4434f3a
    Using config file: /home/me/.signr/config.yaml
    signing on message: signr_0_SHA256_SCHNORR_arbitrary_protocol_string_a61e297a6a401ceef85dff780e757d55016e402230a8c5e308fcb13cf4434f3a
    sig1p4qq4z47zupuxenmn34lfa45erjfwjxnxcns7ll3yz534xlmu4l9gg7c0pfy0vfy3xdez42cuaj7h2q9805a82lxs2gk49reday90wgvelhcv

### verify



## features

It provides the following functionality:

- [x] Key generation - using the system's strong entropy source
    - [x] hexadecimal
    - [x] nostr nsec format
- [x] Secret key import
    - [x] hexadecimal
    - [x] nostr nsec key format
- [x] Signing of data
    - [x] from file
    - [x] piped via stdin
    - [ ] on raw hash (v1.1.0)
- [x] Verification - checking a signature matches
    - [x] from file
    - [x] piped via stdin
- [x] Keychain management
    - [x] storing keys in user profile
    - [x] validating filesystem security of keychain files and folder
    - [x] setting a default key to use when unspecified for signing
    - [x] Encryption of private keys.

In order to prevent cross-protocol attacks, the signature is applied not
directly on the hash of the message, but rather a distinctive structure that
prevents collisions between a signature for this versus any other use of the
same hash functions and elliptic curves.

## Signing Material

For want of a better name, `signr` does not sign directly on hashes of messages, but always a construction that includes a hex encoding of the hash, thus this is an intermediate step to producing a signature, as this "signing material" is then hashed to produce what is signed, and must be reconstructed if the protocol isolates the signature from the rest of the signing material.

The raw bytes that are hashed using SHA256 are constructed as follows:

1. Magic - the string `signr`, the base namespace field.
2. Version - Monotonic number as a string encoding the version being used.
   Starts with 0.
3. The hash function used, `SHA256` normally but allowing future additions, This
   is also part of the signature prefix.
4. The signature algorithm is here encoded. The value `SCHNORR` represents supports BIP-340 style 32 byte
   public keys as specified in NIP-19, and produces 64 byte Schnorr signatures.
   It is specified to enable future expansion to support other signature
   algorithms.
5. Nonce - optional - a strong random value of 64 bits as 16 hex characters, that are
   repeated in the signature prefix to enable the generation of the actual
   message hash that is signed. This is here to ensure the same hash is never
   signed twice as this weakens the security of the EC keys.
    This field is optional as the protocol may not be liable to message hash repetition, such as merkle trees and data that compounds progressively such as DAGs.
6. Custom Protocol String - optional - applications can add a fixed string as a namespace for their signatures in this position.
7. Public Key of the signatory, required for verification.
8. Hash of the message being signed in hex

The string is interpreted as standard ASCII, and the hash that is signed is
generated from these ASCII bytes. All hexadecimal digits are lower case. Hash
function strings should be as they are normally written as identifiers, usually
all caps, this is necessary otherwise there could be malleability. Each section
is separated by a underscore, so the whole string is selected "by word" with
most GUI text selection systems using a double click.

The provided signature contains the Bech32 encoded signature bytes in the place
of the message hash. The verifier splits off this signature, adds the message
hash in its place, hashes the resulting string, and then after decoding the
signature to bytes, calls the secp256k1 schnorr signature verify function and
gets the result. If the signature was made on the hash of this same string then
it will pass.

## Signatures

The base form of signature generated by `signr` essentially takes the format described above and replaces the final hexadecimal hash with a Bech32 signature with the Human Readable Part `nsig` to be consistent with the `nsec` and `npub` HRPs used with Nostr Bech32 keys.

For various reasons, a protocol may omit this prefixing, even it may only store the raw 64 byte signature in binary form.

However, to verify the signature it must have the first 4 fields as described above, and the `npub` public key appended to the hex encoded form of the file's hash (SHA256 presumably) - a detail that will be baked into the protocol elsewhere.

It can also have a custom protocol string, constructed as spaces replaced by hyphens, and all leading and trailing whitespace characters trimmed off, including underscores (this is automatic canonicalization for the custom protocol field). Again, this is a detail that will be already baked into the protocol that consumes this signature library or CLI tooling.

This "signing material" is gathered in order to reconstitute the original hash that will be used with the public key and the verification function to determine if the signature is valid.