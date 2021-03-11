all: stork

stork: _build
	@go build -ldflags="-w -s" -o _build/stork cmd/stork/*.go

_build:
	@mkdir -p _build

clean:
	@rm -rf _build