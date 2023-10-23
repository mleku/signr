# Protocol Buffers

to regenerate the Go code for the protocols, you need to install a few things:

## Protocol Buffers Compiler:

on arch, manjaro et al:

    sudo pacman -S protobuf

on debian, ubuntu, pop OS:

    sudo apt install -y protobuf-compiler
    
check the version:
    
    protoc --version  

    libprotoc 24.3

note that this number changed its meaning after some way into version 3, the major number was dropped and the minor version continued.

## Protocol Buffers Generators for Go

the protocul buffers generator plugins were installed at the time of writing thus:

    # go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

which at time of writing outputs:

    go: downloading google.golang.org/protobuf v1.31.0

then you need the grpc generators, at writing using the @latest gives this output, which are the versions that were run to generate the go files in this package:

    # go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

    go: downloading google.golang.org/grpc v1.59.0

    go: downloading google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0

    go: downloading google.golang.org/protobuf v1.28.1

Newer or older versions of the plugins may produce different results. The output shown above can be used to specify the versions, that can be installed explicitly as follows:

    # go install google.golang.org/protobuf@v1.31.0 ;
      go install google.golang.org/grpc@v1.59.0 ;
      go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 ;
      go install google.golang.org/protobuf@v1.28.1 ;

you will then see where these were placed by the go tool:

    # ls -nl $GOPATH/pkg/mod/google.golang.org
    total 28
    drwxr-xr-x  3 1000 1000 4096 Oct 23 06:52 genproto
    dr-xr-xr-x  7 1000 1000 4096 Oct 23 06:52 genproto@v0.0.0-20230822172742-b8732ec3820d
    drwxr-xr-x  3 1000 1000 4096 Oct 23 06:52 grpc
    dr-xr-xr-x 24 1000 1000 4096 Oct 23 06:52 grpc@v1.3.0
    dr-xr-xr-x 35 1000 1000 4096 Oct 23 06:52 grpc@v1.59.0
    dr-xr-xr-x 12 1000 1000 4096 Oct 23 06:52 protobuf@v1.28.1
    dr-xr-xr-x 13 1000 1000 4096 Oct 23 06:52 protobuf@v1.31.0
