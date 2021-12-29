package wav

import "testing"

func TestInt16Uint16(t *testing.T) {
	var u uint16 = 32767 + 1
	act, exp := int16(u), int16(-32768)
	if act != exp {
		t.Errorf("Incorrect convertion uint16->int16: want %v, got %v", exp, act)
	}
}
