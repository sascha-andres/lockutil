#!/usr/bin/env fish

go build -o lock_test main.go
./lock_test -verbose -port 51001
./lock_test list -port 51001
./lock_test force-release -verbose -port 51001 -force-token abc
./lock_test list -port 51001
rm lock_test
