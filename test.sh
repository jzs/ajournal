#! /bin/bash
echo 'Running linter'
golint $(go list ./... | grep -v vendor | grep -v .git)
echo 'Running go vet'
go vet $(go list ./... | grep -v vendor | grep -v .git)
# echo 'Running unit tests'
# go test -cover=true $(go list ./... | grep -v vendor)
echo 'Running unit and integration tests'
go test -cover=true -tags=integration $(go list ./... | grep -v vendor | grep -v .git)
