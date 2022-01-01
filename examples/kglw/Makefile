.PHONY: all clean

all: play

play: kglw.wav
	mplayer $^

kglw.wav: *.go
	go run .

clean:
	rm -f kglw.wav


