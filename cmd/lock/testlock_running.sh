#!/usr/bin/env fish

./lock -verbose
./lock -verbose -timeout 2

sleep 2

./lock release -verbose
