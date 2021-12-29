.PHONY: all clean

all: play

play: sample.wav
	mplayer $^

sample.wav: *.go
	go run .

clean:
	rm -f sample.wav
