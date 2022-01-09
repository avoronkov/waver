.PHONY: all run clean midisynth controller

all: run

run: midisynth
	./midisynth

midisynth:
	PKG_CONFIG_PATH=/usr/lib/pkgconfig go build ./cmd/midisynth

controller:
	PKG_CONFIG_PATH=/usr/lib/pkgconfig go build ./cmd/controller

clean:
	rm -f ./midisynth
