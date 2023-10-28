BINARY_NAME=dnscheckup

build:
	env GOOS=linux go build -o bin/linux/${BINARY_NAME} main.go
	env GOOS=windows go build -o bin/windows/${BINARY_NAME}.exe main.go

run:
	go run main.go

clean:
	go clean
	rm -f bin/linux/${BINARY_NAME}
	rm -f bin/windows/${BINARY_NAME}.exe