# Being on Linux

Compile for Linux:
```
go build -o engine.out engine.go
```

Compile for Windows:
```
GOOS=windows GOARCH=amd64 go build -o engine.exe engine.go
```