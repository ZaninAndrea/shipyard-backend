env GOOS=windows go build -o build/generic_server.exe
env GOOS=darwin go build -o build/generic_server_macos
env GOOS=linux go build -o build/generic_server_linux