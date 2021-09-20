all:
	go build
	env GOOS=linux GOARCH=amd64 go build -o 1bin.x64
	env GOOS=windows GOARCH=amd64 go build -o 1bin.exe
	env GOOS=linux GOARCH=arm GOARM=7 go build -o 1bin.arm7
	env GOOS=linux GOARCH=arm64 go build -o 1bin.aarch64
