ghmd
====

Simple Go CLI tool that renders GitHub Markdown files using the GitHub API.

Compatibility
-------------

It should work on any system supported by Go. For now, I only tested on Linux but if you use it on another system, feel free to update this README and send me a pull request :)

Build
-----

The only required dependency is Go. Although it should work with any version of Go, I only have tested with Go1.2.

You can get and build `ghmd` by issuing the following command:

```
go get github.com/gilliek/ghmd
```

(make sure your `$GOPATH` is correctly set)

Download binary
---------------

If you want to download a binary file, you can use [Gobuild.IO](http://gobuild.io/download/github.com/gilliek/ghmd)

Usage
-----

Assuming that `ghmd` is in your `$PATH`:

```
ghmd README.md
```

You can also use the `-w` switch to automatically update the generated HTML file when the `README.md` is modified:

```
ghmd -w README.md
```

For a full description of the options, please use:

```
ghmd -h
```

