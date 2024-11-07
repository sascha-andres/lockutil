#!/usr/bin/env fish

go build -o lock_test main.go

./lock_test -verbose -port 51001
./lock_test -verbose -timeout 2 -port 51001

./lock_test list -port 51001

./lock_test release -verbose -port 51001

./lock_test list -port 51001

rm lock_test
