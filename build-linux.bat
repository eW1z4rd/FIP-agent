rmdir .\build\ /s /q
mkdir .\build\
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ./build/fip_agent_linux_amd64 main.go
pause