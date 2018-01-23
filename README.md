# backoff

# Install

```sh
$ go get -u github.com/josedonizetti/backoff/cmd/...
$ backoff --help
```

# Example
```sh
$ backoff https://httpbin.org/delay/3
$ backoff -a2 https://httpbin.org/delay/3
```

# Docker
```sh
$ docker run josedonizetti/backoff https://httpbin.org/delay/3
```

# Design Decisions

- [Structured Log](https://peter.bourgon.org/go-best-practices-2016/#logging-and-instrumentation)
- [Logger as a dependecy](https://peter.bourgon.org/go-best-practices-2016/#top-tip-10)
- [kingpin](https://github.com/alecthomas/kingpin) for command-line flags
- [godep](https://github.com/tools/godep) dependecy management
