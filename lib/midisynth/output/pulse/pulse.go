package pulse

/*
#cgo pkg-config: libpulse-simple
#include <pulse/simple.h>
#include <pulse/error.h>
*/
import "C"
import (
	"fmt"
	"io"
	"time"
	"unsafe"
)

type Pulse struct {
	Spec C.pa_sample_spec
	Pa   *C.pa_simple
}

func New(sampleRate, channels, bits int) (*Pulse, error) {
	var format C.pa_sample_format_t
	switch bits {
	case 2:
		format = C.PA_SAMPLE_S16LE
	default:
		panic(fmt.Errorf("Unsupported bits: %v", bits))
	}
	spec := C.pa_sample_spec{
		format:   format,
		rate:     C.uint(sampleRate),
		channels: C.uchar(channels),
	}
	var cerr C.int
	s := C.pa_simple_new(nil, C.CString("go-pulse"), C.PA_STREAM_PLAYBACK, nil, C.CString("playback"), &spec, nil, nil, &cerr)
	if s == nil {
		return nil, fmt.Errorf("pa_simple_new failed: %v", C.GoString(C.pa_strerror(cerr)))
	}
	return &Pulse{
		Spec: spec,
		Pa:   s,
	}, nil
}

func (p *Pulse) Close() error {
	var cerr C.int
	if C.pa_simple_drain(p.Pa, &cerr) < 0 {
		return fmt.Errorf("pa_simple_drain failed: %v", C.pa_strerror(cerr))

	}
	C.pa_simple_free(p.Pa)
	return nil
}

func (p *Pulse) Play(r io.Reader) error {
	// read 5ms at once
	buf := make([]byte, int(p.Spec.rate)*int(p.Spec.channels)*2/100)
	var cerr C.int
	// var frames int64
	playStart := time.Now()
	var writtenMicro time.Duration
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// slog.Info("n", "n", n)
		if C.pa_simple_write(p.Pa, unsafe.Pointer(&buf[0]), C.ulong(n), &cerr) < 0 {
			return fmt.Errorf("pa_simple_write failed: %v", C.GoString(C.pa_strerror(cerr)))
		}

		writtenMicro += 10 * time.Millisecond
		passed := time.Since(playStart)
		ahead := writtenMicro - passed
		// slog.Info("written", "writtenTime", writtenMicro, "passedTime", passed, "ahead", ahead)
		if ahead > 5*time.Millisecond {
			// slog.Info("sleep", "d", ahead-5*time.Millisecond)
			time.Sleep(ahead - 5*time.Millisecond)
		}
	}
	return nil
}
