install:
	go run make.go -install

crosscompile:
	go run make.go -crosscompile

clean:
	rm -rf pkg src
