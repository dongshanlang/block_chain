NAME=blockChain
DB=blockChain.db
GOOS = darwin
GOARCH = amd64


build: clean
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build  -o $(NAME) ./*.go

clean:
	rm -rf ./$(NAME)
	rm -rf ./$(DB)


