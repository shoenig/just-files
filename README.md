# just-files

![GitHub](https://img.shields.io/github/license/shoenig/just-files?style=flat-square)

`just-files` is a trivial static files server in Go.

### Features

Basically none

### Usage

Set `BIND` environment variable to configure bind address

Set `PORT` environment variable to configure bind port

Arguments are a comma-separated list of `<url-path>:<file-path>` pairs.

### Examples

```shell
BIND=0.0.0.0 PORT=8081 just-files /:/www
```

### Building

The `just-files` file server is written in Go. IT can be built using the normal Go toolchain steps, e.g.

```shell
go build
```

### Installing

With the Go toolchain installed, 

```shell
go install github.com/shoenig/just-files@latest
```

### Contributing

Open an issue

### License

The `just-files` file server is made open source under the [MPL-2.0](LICENSE) license.
