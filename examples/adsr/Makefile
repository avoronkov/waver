.PHONY: all clean

all: play

play: signal.wav
	mplayer $^

signal.wav: *.go
	go run .

clean:
	rm -f signal.wav

