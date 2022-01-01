.PHONY: all clean midisynth

all: midisynth

midisynth:
	PKG_CONFIG_PATH=/usr/lib/pkgconfig go build ./cmd/midisynth
	./midisynth

clean:
	rm -f ./midisynth
