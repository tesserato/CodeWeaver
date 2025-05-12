# go test -cover -v

go test -coverprofile="coverage.out"
go tool cover -func="coverage.out"
go tool cover -html="coverage.out"