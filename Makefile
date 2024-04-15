clean:
	go clean

compile: 
	@make clean
	if [ ! -d "./bin" ]; then mkdir ./bin; fi
	@go build -o ./bin/stewel ./cmd/stewel/main.go

compile_linux: 
	@make clean
	if [ ! -d "./bin" ]; then mkdir ./bin; fi
	GOARCH=amd64 GOOS=linux go build -o ./bin/stewel-linux-amd64 ./cmd/stewel/main.go

compile_darwin: 	
	@make clean
	if [ ! -d "./bin" ]; then mkdir ./bin; fi
	GOARCH=arm64 GOOS=darwin go build -o ./bin/stewel-darwin-arm64 ./cmd/stewel/main.go

run: 
	@make clean
	@make compile
	./bin/stewel

dev:
	@go run cmd/stewel/main.go