Stork is a small utility that aims to automate and simplify some tasks related to software release cycles such as 
reading the current version from a file, prompt for a new version, create and push git tags and so on.

## Install

    # make sure go modules are used
    GO111MODULE=on go get github.com/evilsocket/stork/cmd/stork

You can run a file with `stork -f /path/to/file.stork`, use `stork -h` for a list of all the options.

## Example

    #!/usr/bin/env stork -f
    # parse the current version from this file
    version:file "example_version.go"
    # ask for a new version and update the file
    version:from_user
    # with the new version defined, commit if needed and create a new git tag
    git:create_tag $VERSION
    # build the docker image and tag it as 'latest'
    docker:build "example/project", ".", "latest"
    # tag 'latest' with the current version and push it to docker hub
    docker:create_tag "example/project", $VERSION, "latest"

## Commands

### Shell

`shell:do "<COMMAND>`

Execute a command with the current `$SHELL`.

### Version

`version:file "<FILE NAME>"`

Read the current version from the specified file, sets `$VERSION` and `$VERSION_FILE`.

`version:read "<FILE NAME>", "<VAR NAME>"`

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
