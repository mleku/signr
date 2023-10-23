# signr

designed along similar lines to `ssh-keygen`, `signr` can generate new keys,
password protect them, keeping them in a directory with a familiar similar
layout as `.ssh`, as well as sign and verify.

## installing

if your go installation is all set up, and you have done all the things, you should be able to install `signr` as so:

    go install github.com/mleku/signr

### setting up your Go dev environment

if you don't know what all the things are, then read on:

#### download Go

download a recent version of Go. as at writing, that's v1.20.10:

find them here: https://go.dev/dl/

no idea why they are calling 1.21 branch stable, since it's an odd numbered minor. they are WRONG, so get 1.20:

    cd $HOME 

note that `cd` alone or `cd ~` both have the same result in most linux and posix compliant systems (YMMV with windows 'posix' compliance.)

then

    wget https://go.dev/dl/go1.20.10.linux-amd64.tar.gz

i run from the source from the source repository of the original go source code (not github) but i'm not going to explain how to do that.

mac users, just remember: apple loves you, and since i don't like apple, you are on your own. something something brew something tada. don't forget to install the xcode first of course.

windows, just install WSL and bash terminal, then read WSL where i say 'ubuntu'.

#### install Go

then, unpack the binary distribution as so:

    tar xvf path/to/download/go1.20.10.linux-amd64.tar.gz

#### configure your shell environment

put the following lines at the end of your `~/.bashrc` or `~/.zshrc` or whatever your preferred shell's startup script is:

    export GOBIN=$HOME/.local/bin
    export GOPATH=$HOME
    export GOROOT=$GOPATH/go
    export PATH=$GOBIN:$GOROOT/bin:$PATH

close your current shell session (ctrl-D) and log in again/open up a new session and you should be able to do this:

    go version

which will print something like:

    go version go1.20.10 linux/amd64

#### get the source code

assuming you have installed essentials git on your system... for that:

##### arch linux

    sudo pacman -S --noconfirm git wget curl base-devel

##### debian/ubuntu/pop OS

    sudo apt -y install build-essential git wget curl

then clone the source code:

    git clone https://github.com/mleku/signr.git

or if you have a github account, you can use the SSH link instead:

    git clone git@github.com:mleku/signr.git

#### compile and install

then you can run the following to place the `signr` binary in your $PATH:

    cd signr
    go install .

to build the RPC API, go visit [pkg/protobufs/README.md](pkg/protobufs/README.md)

## usage

### help

The first thing to know is how to access the inbuilt help. The `signr` command by default prints the top level help information with no arguments:

    signr

will print something like this:

    signr - A command line interface for generating, importing, signing, verifying and managing keys used with the Nostr protocol.
    
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

    signr sign go.sum

which will return something like:
    
    signr_0_SHA256_SCHNORR_088b9a5b4bf3fd87_npub1vp0q3vl6sgvpwq5pcfdjyjj8q2kg3r6aa45vclzu8xuyrmaahl4svlvppa_sig135x8pmhfe2ancwjxeeu3rfdumtr95ceaqcw2vx2nv2ek655hyq5unc4g968r8jt4ucflqsz7m2l9jh2dy36va5kr4nzrgpgvrt2gzjcph9zgu

Note the presence of the random nonce and the public key.

For other uses, it may be desired to only sign a hash and only get the signature back.

If the first parameter is 64 characters long, and purely containing `0-9a-f` it will be interpreted as a hash and signed without the addition of the random number as with the `filename/-` and will return only the signature in bech32.

    signr sign --hex 5167fc037d96be70b809e3c669b7cb64864206273757f670af8f729beabf6ddc

which will return something like:

    type password to unlock encrypted secret key:

    74ca3325925ef0de4006d80cf64aed379596a397850f4514f22baf0f30894149fa3d16d6fa7524271eeb88010d08febc699127a615c409881b7713ebd45c2605

The signature will be printed in Bech32 by default, use the flag `--hex` to get a hexadecimal version:

    signr sign --hex 5167fc037d96be70b809e3c669b7cb64864206273757f670af8f729beabf6ddc

which will return something like:
    
    033eb0e02fd0caf521db27d84825a0bd691c29332f4a298de4a499886fd2f52139600625d833ceb7bb9b876b267b8e85567f3db74d8f7243e0e78ce72f269642


For reasons of prevention of namespace attacks, these signatures are generated on a structured value that is hashed, like so:

    signr_0_SHA256_SCHNORR_5167fc037d96be70b809e3c669b7cb64864206273757f670af8f729beabf6ddc

rather than being signed directly on the provided value encoded as binary.

To allow custom protocols to be devised, this can have an additional custom string using the `--custom` flag:

    signr -v sign --custom arbitrary-protocol-string go.sum

which will return something like:
    
    Using config file: /home/me/.signr/config.yaml
    signing on message: signr_0_SHA256_SCHNORR_arbitrary-protocol-string_8b716747195afc47_npub1vp0q3vl6sgvpwq5pcfdjyjj8q2kg3r6aa45vclzu8xuyrmaahl4svlvppa_a61e297a6a401ceef85dff780e757d55016e402230a8c5e308fcb13cf4434f3a
    
    signr_0_SHA256_SCHNORR_arbitrary-protocol-string_8b716747195afc47_npub1vp0q3vl6sgvpwq5pcfdjyjj8q2kg3r6aa45vclzu8xuyrmaahl4svlvppa_sig1wk3ndpaldxct9v9dc3p0r0ynpwzcy9nyuadyzfp7xx7hcrxrm6c4ptanwlqywnf0vdlyc28krnfszyw0apev2t5vyrwetgp6deqvc7ggkf4xy

The arbitrary custom string has any leading or following spaces or carriage returns removed:

    signr -v sign --custom "arbitrary protocol string 
        " go.sum

which will return something like:
    
    Using config file: /home/me/.signr/config.yaml
    signing on message: signr_0_SHA256_SCHNORR__41a7170e39e1a703_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52
    type password to unlock encrypted secret key:
    secret decrypted: true; decrypted->pub: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs, stored pub; npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs
    
    signr_0_SHA256_SCHNORR__41a7170e39e1a703_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_nsig13qu78yepuzl7z6zafj388h34funs7e2aku5naue8a0q7v07gp5y8wprstnk6ggcwgvf0utuvuuqe9dkk9ta6grdc9nh5frlxulkj3asjlr2th

Carriage returns as well, though these can only end up in the text very deliberately with a copy/paste or other programmatic method:

    signr -v sign --custom "arbitrary protocol string 
        " go.sum

which will return something like:
    
    Using config file: /home/me/.signr/config.yaml
    signing on message: signr_0_SHA256_SCHNORR__d73d19ff8de5dcec_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52
    type password to unlock encrypted secret key:
    secret decrypted: true; decrypted->pub: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs, stored pub; npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs

    signr_0_SHA256_SCHNORR__d73d19ff8de5dcec_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_nsig1xucg87qfca6capz2j5mzw9lhkx3vm3j47htqqgp7z5nrvre950zzwr6n3udkwf0lz66lttxxh52n6hzptx6f8he40msngwk90ww9rhcgmfljx

For completeness, the same thing except with a hash instead of a file:

    signr -v sign --custom "arbitrary protocol string 
    " a61e297a6a401ceef85dff780e757d55016e402230a8c5e308fcb13cf4434f3a

which will return something like:
    
    Using config file: /home/me/.signr/config.yaml
    signing on message: signr_0_SHA256_SCHNORR__npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_a61e297a6a401ceef85dff780e757d55016e402230a8c5e308fcb13cf4434f3a
    type password to unlock encrypted secret key:
    secret decrypted: true; decrypted->pub: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs, stored pub; npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs
    
    nsig1hkdkdaulkyky43ec077phk0dlmpwnmw65yvh0vg6w5mcj2xqjnw55de9cftl45gjdts32e2j7yjp968v6l7wfutgefzx03yay79af5ssy52q0

### verify

The `verify` command with a filename or piped file input, plus a signature, or signature filename enables the signing of a file.

A simple all-in-one example that covers a lot of the features involved, including a signing using an encrypted key with an environment variable, looks like this:

    signr -v -c verify go.sum `SIGNR_PASS=aoeu signr -v -c sign go.sum`

which will return something like:
    
    > Using config file: /home/me/.signr/config.yaml
    > signing on message: signr_0_SHA256_SCHNORR_503ca48a2056127e_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52
    > secret decrypted: true; decrypted->pub: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs, stored pub; npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs
    > Using config file: /home/me/.signr/config.yaml
    > pubkey input: 
    > pubkey from env: 
    > nonce found 503ca48a2056127e
    > pubkey in signature: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs
    > loading pubkey: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs
    > adding nonce 
    > message: 
    
    signr_0_SHA256_SCHNORR_503ca48a2056127e_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52
    VALID

In this example several things are visible - as the logging has been enabled with `-v`:

The actual string that is hashed to generate a signature is this:

    signr_0_SHA256_SCHNORR_503ca48a2056127e_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52

As you can see later on in the output, this text is repeated in the verification step as this full string is read in, in this case, with a standard `signr` signing material format. 

This example does not show the output of the embedded signing command, which looks like this:

    SIGNR_PASS=aoeu signr -v -c sign go.sum

which will return something like:
    
    > Using config file: /home/me/.signr/config.yaml
    > signing on message: signr_0_SHA256_SCHNORR_08f0f2582204bead_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52
    > secret decrypted: true; decrypted->pub: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs, stored pub; npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs
    
    signr_0_SHA256_SCHNORR_08f0f2582204bead_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_nsig1qvyx4xhs8nek3uecpeleqvqn8nxrvq30zpfxj0ephf7pg6vxh0xrme7wrp9ufwdxv453q0lj22yw5608h5dwypux4esy4yndrqwv48qv6szf5

The actual signature, that could be provided as a file alongside a download, for example, or in the advertisment of a file available on a torrent, let's say, is this:

    signr_0_SHA256_SCHNORR_08f0f2582204bead_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_nsig1qvyx4xhs8nek3uecpeleqvqn8nxrvq30zpfxj0ephf7pg6vxh0xrme7wrp9ufwdxv453q0lj22yw5608h5dwypux4esy4yndrqwv48qv6szf5

### advanced usage

If the protocol separates things, then first of all, the nonce will not be used. The protocol would need to use the `--custom` feature to embed any extra strings, which are sanitised to be only printable characters with single spaces represented as hyphens.

When a raw hexadecimal hash is provided in place of a file, `signr` automatically changes its mode and will output either an `nsig` formatted signature alone, or a hexadecimal signature, as requested. These two look like this:

    SIGNR_PASS=aoeu signr -v -c sign --sigonly go.sum

which will return something like:
    
    > Using config file: /home/me/.signr/config.yaml
    > signing on message: signr_0_SHA256_SCHNORR_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52
    > secret decrypted: true; decrypted->pub: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs, stored pub; npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs
 
    nsig1sp0u578eds7mka2v4gxrrcyyzmpkm9vkrdp6ddyvr7sjrlxprjrrcd6m6y2cnq2c6pf2ze7utc4y877x94amludluxja52dvk9z88pg2hzgyr

As you can see, the final output, which has been separated visually (in the commandline the lines starting with `>` are also in a different color) is just the signature. Changing it to `--hex` instead of `--sigonly` yields much the same output but with a different final result:

    SIGNR_PASS=aoeu signr -v -c sign --hex go.sum

which will return something like:
    
    > Using config file: /home/me/.signr/config.yaml
    > signing on message: signr_0_SHA256_SCHNORR_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52
    > secret decrypted: true; decrypted->pub: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs, stored pub; npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs
    
    805fca78f96c3dbb754caa0c31e08416c36d95961b43a6b48c1fa121fcc11c863c375bd115898158d052a167dc5e2a43fbc62d7bbff1bfe1a5da29acb1447385

Again the input and output lines have been visually spaced as they are colorised in the actual output using the `-c` (`--color`) flag, but both are the same value once parsed and decoded into binary.

Lastly, to illustrate what happens when a hash is provided instead of a file, the result is the same as these `--sigonly` and `--hex` versions but you don't specify a file, instead you just give it the hash in hexadecimal:

    SIGNR_PASS=aoeu signr -v -c sign --hex 830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52

which will return something like:
    
    > Using config file: /home/me/.signr/config.yaml
    > signing on message: signr_0_SHA256_SCHNORR_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52
    > secret decrypted: true; decrypted->pub: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs, stored pub; npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs
    
    805fca78f96c3dbb754caa0c31e08416c36d95961b43a6b48c1fa121fcc11c863c375bd115898158d052a167dc5e2a43fbc62d7bbff1bfe1a5da29acb1447385
    
In the above the `--hex` flag is used, and we get the same raw hex as you saw above from the file (the example uses the current state of the `go.sum` file from the repository). If you explicitly use `--sigonly` or no flag at all, if the input is a hex (and again, for brevity, we are signing with the default signature, providing its password via environment variable), you get the `nsig` instead:

    SIGNR_PASS=aoeu signr -v -c sign 830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52

which will return something like:
    
    > Using config file: /home/me/.signr/config.yaml
    > signing on message: signr_0_SHA256_SCHNORR_npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs_830bcbdbcf0b55307030e0838752af4e79fecca4ee0d27602f8e6cba6239bd52
    > secret decrypted: true; decrypted->pub: npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs, stored pub; npub1v79pv39a3asvwrcwct6qs6jpupr2uk707mz50c0ea7yrtg4jl32sxhhmvs
    
    nsig1sp0u578eds7mka2v4gxrrcyyzmpkm9vkrdp6ddyvr7sjrlxprjrrcd6m6y2cnq2c6pf2ze7utc4y877x94amludluxja52dvk9z88pg2hzgyr



## features

It provides the following functionality:

- [x] Key generation - using the system's strong entropy source
    - [x] hexadecimal secret key
    - [x] nostr nsec format public key
- [x] Secret key import
    - [x] hexadecimal secret key
    - [x] nostr nsec secret key format
- [x] Signing of data
    - [x] from file
    - [x] piped via stdin
    - [x] on raw hash
    - [x] arbitrary custom namespace field
- [x] Verification - checking a signature matches
    - [x] from file
    - [x] piped via stdin
    - [x] external additional pubkey and custom field for signature only
- [x] Keychain management
    - [x] storing keys in user profile
    - [x] validating filesystem security of keychain files and folder
    - [x] setting a default key to use when unspecified for signing
    - [x] encryption of private keys.
- [ ] gRPC/Protobuf API

In order to prevent cross-protocol attacks, the signature is applied not
directly on the hash of the message, but rather a distinctive structure that
prevents collisions between a signature for this versus any other use of the
same hash functions and elliptic curves.

## Signing Material

For want of a better name, `signr` does not sign directly on hashes of messages, but on "signing material", always a construction that includes a hex encoding of the hash, thus this is an intermediate step to producing a signature, as this **signing material** is then hashed to produce what is signed, and must be reconstructed if the protocol isolates the signature from the rest of the signing material.

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