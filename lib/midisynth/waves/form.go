package waves

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Form struct {
	values  []float64
	nValues int
}

var _ Wave = (*Form)(nil)

func ParseFormFile(path string) (*Form, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseFormInput(f)
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func parseFormInput(in io.Reader) (*Form, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	intValues := []int{}
	maxVal := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		val, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		intValues = append(intValues, val)
		if a := abs(val); a > maxVal {
			maxVal = a
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(intValues) == 0 {
		return nil, fmt.Errorf("Empty list of Form values")
	}

	maxFloatValue := float64(maxVal)
	values := make([]float64, len(intValues))
	for i, v := range intValues {
		values[i] = float64(v) / maxFloatValue
	}

	return &Form{
		values:  values,
		nValues: len(values),
	}, nil
}

func (f *Form) Value(tm float64, ctx *NoteCtx) float64 {
	t := tm - ctx.Period*float64(int(tm/ctx.Period))
	tf := t * ctx.Freq * float64(f.nValues)
	idx := int(t * ctx.Freq * float64(f.nValues))
	dtf := tf - float64(idx)
	y1 := f.values[idx]
	y2 := f.values[(idx+1)%f.nValues]
	y := y1 + dtf*(y2-y1)
	return y
}
