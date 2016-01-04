install:
	go get github.com/mitchellh/gox
	gox -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}"

clean:
	rm -rf pkg src
