.PHONY: all run clean midisynth install

all: midisynth

run: midisynth
	./midisynth

midisynth:
	PKG_CONFIG_PATH=/usr/lib/pkgconfig go build ./cmd/midisynth

install:
	PKG_CONFIG_PATH=/usr/lib/pkgconfig go install ./cmd/midisynth

clean:
	rm -f ./midisynth
