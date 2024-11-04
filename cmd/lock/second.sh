#!/usr/bin/env fish

go build -o lock_test2 main.go

./lock_test2 -timeout 10 -verbose

echo a

./lock_test2 release -verbose

rm lock_test2
