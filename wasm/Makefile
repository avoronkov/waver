.PHONY: all clean server static

TARGETS= static server

all: $(TARGETS)

static:
	env GOOS=js GOARCH=wasm go build -o static/assets/demo.wasm ./cmd/wasm
	gzip -f -9 static/assets/demo.wasm
	cd src && yarn install
	cp -rf ./src/node_modules/codemirror/lib ./src/node_modules/codemirror/addon ./static/assets

server:
	go build -o $@ ./cmd/server

clean:
	rm -f server
