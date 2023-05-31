#!/usr/bin/env just --justfile

update:
  go get -u
  go mod tidy -v

install:
  cd cmd/ask && go install .

gif:
    vhs .github/ask.tape -o .github/ask.gif

# release, e.g `just release v0.12`
release version:
    sed -i 's/Version = ".*"/Version = "{{version}}"/' ./cmd/ask/const.go
    git add .
    git commit -m "release {{version}}"
    git push
    git tag {{version}}
    git push --tags
    rm -rf dist
    goreleaser release
