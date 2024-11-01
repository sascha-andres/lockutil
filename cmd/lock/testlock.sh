#!/usr/bin/env fish

go build -o lock_test main.go

./lock_test -verbose
./lock_test -verbose -timeout 2

sleep 2

./lock_test release -verbose

rm lock_test
