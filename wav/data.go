package wav

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Data struct {
	Samples []int16
}

func (d *Data) String() string {
	return fmt.Sprintf("Data: bytes=%v, samples=%v", d.ChunkSize(), len(d.Samples))
}

func ParseWavData(buf []byte) *Data {
	data := &Data{}
	l := len(buf) / 2
	data.Samples = make([]int16, l)

	for i := 0; i < l; i++ {
		data.Samples[i] = int16(binary.LittleEndian.Uint16(buf[i*2:]))
	}
	return data
}

func (d *Data) AddSample(value int16) {
	d.Samples = append(d.Samples, value)
}

func (d *Data) AddStereoSample(value int16) {
	d.Samples = append(d.Samples, value, -value)
}

func (d *Data) FullSize() uint32 {
	return uint32(len("data")) + 4 + d.ChunkSize()
}

func (d *Data) ChunkSize() uint32 {
	return uint32(len(d.Samples)) * 2
}

func (d *Data) Write(w io.Writer) error {
	io.WriteString(w, "data")
	binary.Write(w, binary.LittleEndian, d.ChunkSize())
	for _, sample := range d.Samples {
		binary.Write(w, binary.LittleEndian, sample)
	}
	return nil
}
