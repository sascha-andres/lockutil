# lockutil

Ever wanted to coordinate shell scripts? Use lockutil to acquire a lock or fail:

Script 1:

```
lock
trap "lock -free" INT EXIT

sleep 1000
```

Script 2:

```
lock
trap "lock -free" INT EXIT

echo "a"
```

Will print 'a' only after script 1 has finished.

## Options

### -timeout
Wait for this number of seconds to acquire lock. If it takes longer fail.
### - name
Use lock name instead of default lock
