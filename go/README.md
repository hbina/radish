# Introduction

Redis implementation in Go.

The aim of this project is to implement Redis in idiomatic-ish Go.
It will try to be as stupid and simple as possible.
Optimizations will come later.
It

# Get involved!

This project is in _work-in-progress_, so share ideas, code and have fun.

The goal is to have all features and commands like the actual [redis](https://github.com/redis/redis) written in C have.
I am always open to collaborate!

# Test

## Running Go tests

There are some tests written in Go.
I aim to make this comprehensive in the future.

```bash
$ go run cmd/main.go <PORT>
$ PORT=<PORT> go test
```

## Running Reference Redis' Tests Suite

To check agains the reference implementation, you can also run the test from `redis` and point it to this implementation.

Note that there are some tests that are irrevelant to us because it is an implementation details (most tests that uses `OBJECT` command) of the reference `redis`.
I maintain a [branch](https://github.com/hbina/redis/tree/hbina-retrofitting-tests-for-go-redis) in my fork of `redis` that trims out some of these.

After cloning the `redis` repository,

```
./runtest --host 127.0.0.1 --port 6380 --tags -needs:repl --ignore-encoding
```

From this command, we are currently ignoring replication and encoding features.
Or if you want to only execute a certain test, you can do,

```
./runtest --host 127.0.0.1 --port 6380 --tags -needs:repl --ignore-encoding --single unit/types/set
```

[Link](https://github.com/redis/redis/blob/203b12e41ff7981f0fae5b23819f072d61594813/tests/README.md) for explanations of some these options.

### Test

# Roadmap

- [x] Client connection / request / respond
- [x] RESP protocol
- [x] able to register commands
- [x] in-mem database
- [x] active key expirer
- [ ] Implementing data structures
  - [x] String
  - [x] List
  - [x] Set
  - [x] Sorted Set
  - [ ] Hash
  - [ ] ...
- [ ] Tests
  - [x] unit/type/set
  - [x] unit/type/string
  - [x] unit/printver

### TODO beside Roadmap

- [ ] Persistence
- [ ] Redis config
  - [ ] Default redis config format
  - [ ] YAML support
  - [ ] Json support
- [ ] Pub/Sub
- [ ] Redis modules
- [ ] Benchmarks
- [ ] master slaves
- [ ] cluster
- [ ] ...
