package dumper

import (
	"encoding/json"
	"log"
	"os"

	"github.com/avoronkov/waver/lib/midisynth/signals"
)

type Json struct {
	file    *os.File
	encoder *json.Encoder
}

var _ signals.Output = (*Json)(nil)

func NewJson(filename string) (*Json, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Json{
		file:    f,
		encoder: json.NewEncoder(f),
	}, nil
}

func (j *Json) ProcessAsync(tm float64, sig signals.Interface) {
	sjson := &SignalJson{
		T:   tm,
		Sig: sig,
	}
	if err := j.encoder.Encode(sjson); err != nil {
		log.Printf("[JSON] Cannot dump %v: %v", sjson, err)
	}
}

func (j *Json) Close() error {
	if j.file != nil {
		return j.file.Close()
	}
	return nil
}

type SignalJson struct {
	T   float64 `json:"T"`
	Sig signals.Interface
}
