set shell := ["bash", "-u", "-c"]

# show just command targets
[private]
default:
    @just --list

# run go test with race detector
[group('build')]
test:
    go test -count=1 -v -race ./...

# scan source with go vet tool
[group('build')]
vet:
    go vet ./...

# tidy go modules
[group('build')]
tidy:
    go mod tidy

# publish latest git tag as a github release
[group('publish')]
release:
        envy exec gh-release goreleaser release --clean

