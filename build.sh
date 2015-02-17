#!/bin/bash

./bin/flatc_mac -g -o ./tmp/ ./scheme.fbs 
cp -f ./tmp/ftunnel/Packet.go src/packet.go
rm -rf ./tmp
