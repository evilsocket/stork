Stork is a small utility that aims to automate and simplify some tasks related to software release cycles such as 
reading the current version from a file, prompt for a new version, create and push git tags and so on.

## Install

    # make sure go modules are used
    GO111MODULE=on go get github.com/evilsocket/stork/cmd/stork

You can run a file with `stork -f /path/to/file.stork`, use `stork -h` for a list of all the options.

## Example

This stork script will parse the current version from example_version.go, then ask the user for a new version and update
this file. It will then push the changes to git and create a new tag with the specified version. The last two lines 
will build, tag and push the docker image for the project.

    #!/usr/bin/env stork -f
    version:file "example_version.go"
    version:from_user

    git:create_tag $VERSION

    docker:build "example/project", ".", "latest"
    docker:create_tag "example/project", $VERSION, "latest"

More in the `examples` folder. For a complete list of commands [refer to the project wiki](https://github.com/evilsocket/stork/wiki/Commands).

## License

Released under the GPL3 license.
