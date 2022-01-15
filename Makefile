.PHONY: all run clean midisynth

all: run

run: midisynth
	./midisynth

midisynth:
	PKG_CONFIG_PATH=/usr/lib/pkgconfig go build ./cmd/midisynth

clean:
	rm -f ./midisynth
