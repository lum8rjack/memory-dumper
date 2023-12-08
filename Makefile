NAME=memory-dumper
BUILD=CGO_ENABLED=0 go build -ldflags "-s -w -buildid=" -trimpath

default: linux

setup:
	go mod init $(NAME)
	go mod tidy

clean:
	rm -f $(NAME)

linux:
	@echo "Compiling for Linux x64"
	GOOS=linux GOARCH=amd64 $(BUILD) -o $(NAME)

arm:
	@echo "Compiling for Linux Arm"
	GOOS=linux GOARCH=arm $(BUILD) -o $(NAME)
