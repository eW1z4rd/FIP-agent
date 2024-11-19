rmdir .\build\ /s /q
mkdir .\build\
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o ./build/fip_agent_windows_amd64.exe main.go
pause