ir-gen-cli-bundler
==================
`ir-gen-cli-bundler` is a Go library that embeds the verification CLI published by the [Conjure repository](https://github.com/palantir/conjure).
The Conjure repository publishes a CLI that transforms Conjure YAML files (either a single file or a directory of files) 
into the JSON intermediate representation (IR). This library fully embeds the CLI published by that repository in Go
source code using the [go-bindata](https://github.com/go-bindata/go-bindata) library and provides a Go API for working
with the CLI.

The Go library works as follows:
* The library fully embeds a specific version of the Conjure CLI in source
* When library functionality that requires the embedded CLI is invoked, the embedded CLI is written to disk and invoked
  * The embedded CLI data is written to `{{tmp}}/_conjureircli/conjure-{{version}}`, where `{{tmp}}` is the directory 
    returned by `os.TempDir()` and `{{version}}` is the version of the CLI embedded in the library
  * If the CLI already exists in that location, it is invoked directly (not written out)
  
Note that, currently, the Conjure CLI is written in Java, and thus invoking the CLI requires the Java runtime.

Updating the bundled CLI
------------------------
To update the version of the CLI bundled in source, do the following:

* Determine the new version of Conjure (it must be available at https://bintray.com/palantir/releases/conjure) 
* Update the value of the `conjureVersion` constant in `conjureircli/generator/generate.go` to the desired version
* Run `./godelw generate` to embed the updated version in source
