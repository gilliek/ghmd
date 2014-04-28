ghmd
====

Simple Go CLI tool that renders GitHub Markdown files using the GitHub API.

### Build

The only required dependency is Go. Although it should work with any version of Go, I only have tested with Go1.2.

Issue the following command to build `ghmd`:

```go build ghmd.go```

### Usage

Assuming that `ghmd` is in your `$PATH`:

```ghmd README.md > output.html```
