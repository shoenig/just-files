set shell := ["bash", "-u", "-c"]

# print just command targets
default:
    @just --list

# run go test with race detector
test:
    go test -count=1 -v -race ./...

# scan source with go vet tool
vet:
    go vet ./...

# tidy go modules
tidy:
    go mod tidy

# publish latest git tag as a github release
release:
        envy exec gh-release goreleaser release --clean

