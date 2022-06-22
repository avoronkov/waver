package wav

import (
	"encoding/binary"
	"io"
)

type DataBytes struct {
	Samples []byte
}

var _ DataInterface = (*DataBytes)(nil)

func (d *DataBytes) Write(w io.Writer) error {
	_, _ = io.WriteString(w, "data")
	_ = binary.Write(w, binary.LittleEndian, d.chunkSize())
	_, err := w.Write(d.Samples)
	return err
}

func (d *DataBytes) chunkSize() uint32 {
	return uint32(len(d.Samples))
}

func (d *DataBytes) FullSize() uint32 {
	return uint32(len("data")) + 4 + d.chunkSize()
}
