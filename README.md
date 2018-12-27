
# linux 

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

# windows 

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w"  -o telegram-robot main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"  -o telegram-robot main.go

# 打包，11M
go build -ldflags "-s -w"  -o telegram-robot main.go

# 压缩，3.4M
upx telegram-robot


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"  -o telegram-robot main.go && upx telegram-robot