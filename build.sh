#!/bin/bash

# install by `port install capnproto
# $ go get -u -t github.com/glycerine/go-capnproto
# $ cd $GOPATH/src/github.com/glycerine/go-capnproto
# $ make # will install capnpc-go and compile the test schema aircraftlib/aircraft.capnp, which is used in the tests.
# $ diff ./capnpc-go/capnpc-go `which capnpc-go`
# $ cp ./capnpc-go/capnpc-go  `which capnpc-go`

capnp compile -ogo src/packet.capnp

go build dtunnel.go

go build stunnel.go