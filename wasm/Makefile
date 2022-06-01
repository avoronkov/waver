.PHONY: all clean server static/assets/demo.wasm

TARGETS= ./static/assets/demo.wasm ./server

all: $(TARGETS)

static/assets/demo.wasm:
	env GOOS=js GOARCH=wasm go build -o $@ ./cmd/wasm

server:
	go build -o $@ ./cmd/server

clean:
	rm -f $(TARGETS)
