#!/usr/bin/env fish

go build -o lock_test main.go

./lock_test -verbose

sleep 5

./lock_test release -verbose

rm lock_test
