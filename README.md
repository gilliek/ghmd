ghmd
====

Simple Go CLI tool that renders GitHub Markdown files using the GitHub API.

Compatibility
-------------

It should work on any system supported by Go. For now, I only tested on Linux but if you use it on another system, feel free to update this README and send me a pull request :)

Build
-----

The only required dependency is Go. Although it should work with any version of Go, I only have tested with Go1.2.

Since it is just a single file without any external dependencies, you can build `ghmd` by simply issuing the following command:

```go build ghmd.go```

Of course, if you are familiar with Go, you can also use the go tools:

```go get github.com/gilliek/ghmd```


Download binary
---------------

If you want to download a binary file, you can use [Gobuild.IO](http://gobuild.io/download/github.com/gilliek/ghmd)

Usage
-----

Assuming that `ghmd` is in your `$PATH`:

```ghmd README.md > output.html```
