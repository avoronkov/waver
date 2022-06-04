.PHONY: all run clean waver install

all: waver

run: waver
	./waver

waver:
	PKG_CONFIG_PATH=/usr/lib/pkgconfig go build ./cmd/waver

install:
	PKG_CONFIG_PATH=/usr/lib/pkgconfig go install ./cmd/waver

clean:
	rm -f ./waver
