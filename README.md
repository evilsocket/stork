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

More in the `examples` folder.

## Commands

### Shell

`shell:do "<COMMAND>"`

Execute a command with the current `$SHELL`.

### Version

`version:file "<FILE NAME>"`

Read the current version from the specified file, sets `$VERSION` and `$VERSION_FILE`.

`version:read "<FILE NAME>", "<VAR NAME>"`

`version:parser "<EXPRESSION>"`

Set the regular expression used by `version:file` and `version:read` to parse the version string. Default to `[Vv]ersion\\s*=\\s*['\"]([\\d\\.ab]+)[\"']`.

Read the version from the specified file and sets `$<VAR_NAME>`.

`version:from_user`

Ask the user for a new version, updates `$VERSION` and `$VERSION_FILE` with the new value.

### Git

`git:create_tag $VERSION` or `git:create_tag "<VERSION>"`

Create and push a new git tag.

### Docker

`docker:repository "<REPOSITORY URL>"`

Set the repository URL for pushing docker images.

`docker:build "<IMAGE NAME>", "<PATH>", "<TAG>"`

Build a docker image from a given path and tag it with the specified tag.

`docker:create_tag "<IMAGE NAME>", "<SOURCE TAG>", "<TARGET TAG>"`

Create and push a new tagged image from a source tag.
