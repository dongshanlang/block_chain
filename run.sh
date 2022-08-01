#!/bin/bash

rm blockChain
rm blockChain.db

go build -o blockChain ./*.go
./blockChain
