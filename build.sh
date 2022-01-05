env GOOS=windows go build -o build/generic_server.exe ./cmd
env GOOS=darwin go build -o build/generic_server_macos ./cmd
env GOOS=linux go build -o build/generic_server_linux ./cmd