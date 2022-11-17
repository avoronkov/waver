.PHONY: all clean server pwa static

TARGETS= pwa

all: $(TARGETS)

pwa:
	go build -o $@ ./cmd/pwa
	GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/pwa

clean:
	rm -f pwa
