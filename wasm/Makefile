.PHONY: all clean server static/assets/demo.wasm.gz

TARGETS= ./static/assets/demo.wasm.gz ./server

all: $(TARGETS)

static/assets/demo.wasm.gz:
	env GOOS=js GOARCH=wasm go build -o static/assets/demo.wasm ./cmd/wasm
	gzip -f -9 static/assets/demo.wasm

server:
	go build -o $@ ./cmd/server

clean:
	rm -f $(TARGETS)
