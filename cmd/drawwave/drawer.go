package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

type Drawer struct {
}

func (d *Drawer) Draw(data []byte, output string) error {
	left, right := d.parseData(data)
	out, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return d.writeData(out, left, right)
}

func (d *Drawer) parseData(data []byte) (left []int16, right []int16) {
	if len(data)%4 != 0 {
		panic(fmt.Errorf("Incorrect buffer size: %v", len(data)))
	}
	l := len(data) / 4
	for i := 0; i < l; i++ {
		s1 := int16(binary.LittleEndian.Uint16(data[i*4:]))
		s2 := int16(binary.LittleEndian.Uint16(data[i*4+2:]))
		left = append(left, s1)
		right = append(right, s2)
	}
	return
}

type TplData struct {
	LeftX, LeftY   string
	RightX, RightY string
}

func NewTplData(left []int16, right []int16) *TplData {
	x := &strings.Builder{}
	for i := range left {
		if i > 0 {
			fmt.Fprintf(x, ",")
		}
		fmt.Fprint(x, i)
	}

	ly := &strings.Builder{}
	for i, y := range left {
		if i > 0 {
			fmt.Fprintf(ly, ",")
		}
		fmt.Fprint(ly, y)
	}

	ry := &strings.Builder{}
	for i, y := range right {
		if i > 0 {
			fmt.Fprintf(ry, ",")
		}
		fmt.Fprint(ry, y)
	}
	return &TplData{
		LeftX:  x.String(),
		LeftY:  ly.String(),
		RightX: x.String(),
		RightY: ry.String(),
	}
}

func (d *Drawer) writeData(w io.Writer, left, right []int16) error {
	tpl, err := template.ParseFiles("./template.html")
	if err != nil {
		return fmt.Errorf("Parsing template failed: %w", err)
	}
	return tpl.Execute(w, NewTplData(left, right))
}
