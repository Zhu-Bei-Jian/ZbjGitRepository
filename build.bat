set GOOS=windows
go build -mod=vendor -o bin/all-in-one.exe ./cmd/all-in-one

SET GOOS=linux
go build -mod=vendor -o bin/all-in-one ./cmd/all-in-one


