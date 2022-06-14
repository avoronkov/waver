package share

import "testing"

func TestEncodeDecode(t *testing.T) {
	input := `% inst sine 'sine'
: 8 -> { sine A4 }`
	encoded, err := Encode(input)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	t.Logf("Encoded data: %v", encoded)

	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded != input {
		t.Errorf("Incorrect decoding: want %q, got %q", input, decoded)
	}
}
