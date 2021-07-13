# redis-dump
- This will help you migrate your data from source to destination redis. Source will always be a single node. Destination can be cluster URL, standalone or twemproxy server.

# usage example:
> ./main copy source:port destination:port --queryString="*"  --scanLimit=1000 --dumpThreads=100 --restoreThreads=100


## Features

- Import data along with TTL
- setup parallelism based in config vars
- option to override keys if already present in destination
- migrate all or specific pattern keys
- it wont delete keys from your source


## Installation & usage
Please clone the repo.
either use exisitng build or build a go binary from source

For mac..
```sh
go build main.go 
```
For linux..
```sh
env GOOS=linux GOARCH=amd64 go build -v main.go 
```

## supported options
| option | description |
| ------ | ------ |
| queryString | prefix pattern to migrate (EX : "*", "USER_*")
| scanLimit | number of keys to be scanned
| dumpThreads | number of threads to get the dump of keys scanned
| restoreThreads| number of threads to restore the keys scanned
| overrideKey | if keys is already present in destination.. override it or not? (defaults to false)


