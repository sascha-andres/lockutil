# lockutil

Ever wanted to coordinate shell scripts? Use lockutil to acquire a lock or fail:

Script 1:

```
lock
trap "lock release" INT EXIT

sleep 1000
```

Script 2:

```
lock -timeout 10
trap "lock release" INT EXIT

echo "a"
```

Will print 'a' only after script 1 has finished.

## lock commands

### no command

acquire lock

### release

release lock

### list

list all active locks

### force-release

force a lock release, a secret token must be provided

## lock options

### -timeout
Wait for this number of seconds to acquire lock. If it takes longer, it fails.

### - port
The port to connect to, defaulting to 50051

### - host
The host to connect to, defaulting to localhost

### - lock
The name of the lock to acquire, defaulting to 'default'

### - help
Prints help a message

### - verbose
Enables verbose logging

## lockd options

### - port
The port to listen on, defaulting to 50051

### - host
The host address to listen on, defaulting to localhost

### - lock
The name of the lock to acquire, defaulting to 'default'

### - help
Prints a help message

### secret-token

pass to enable forcefully unlocks

## as a go package

In `lockutil.go` a client library is provided for use in go applications.