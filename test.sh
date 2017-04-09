echo 'Running linter'
golint $(go list ./... | grep -v vendor)
echo 'Running go vet'
go vet $(go list ./... | grep -v vendor)
echo 'Running unit tests'
go test $(go list ./... | grep -v vendor)
# echo 'Running integration tests'
# go test -tags=integration $(go list ./... | grep -v vendor)
