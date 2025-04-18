default: build

clean:
	rm -rf out

test:
	go test -v ./... -timeout 15s

build:
	mkdir -p out
	go build -o out/painter ./cmd/painter/main.go