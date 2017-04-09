echo 'Running linter'
golint $(go list ./... | grep -v vendor)
echo 'Running unit tests'
go test $(go list ./... | grep -v vendor)
