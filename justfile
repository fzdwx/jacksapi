#!/usr/bin/env just --justfile

update:
  go get -u
  go mod tidy -v

install:
  cd cmd/ask && go install .
