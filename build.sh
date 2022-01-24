env GOOS=windows GOARCH=amd64 go build -o build/generic_server.exe ./cmd
env GOOS=darwin GOARCH=amd64 go build -o build/generic_server_macos ./cmd
env GOOS=linux GOARCH=amd64 go build -o build/generic_server_linux ./cmd