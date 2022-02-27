package common

import (
	"fmt"
)

type KeyNote int64

type KeyNoteT struct {
	Number  int64
	UdpRepr string
}

func (k *KeyNoteT) String() string {
	return k.UdpRepr
}

func ParseStandardNoteCode(code int64) (*KeyNoteT, error) {
	for _, n := range StandardScale {
		if n.Number == code {
			return &n, nil
		}
	}
	return nil, fmt.Errorf("Unknown standard note code: %v", code)
}

func ParseStandardNote(str string) (*KeyNoteT, error) {
	if n, ok := StandardScale[str]; ok {
		return &n, nil
	}
	return nil, fmt.Errorf("Unknown standard note: %v", str)
}

var StandardScale = initStandardScaleKeyNotes()

func initStandardScaleKeyNotes() map[string]KeyNoteT {
	m := map[string]KeyNoteT{}

	var oct int64
	for oct = 1; oct <= 7; oct++ {
		m[fmt.Sprintf("C%d", oct)] = KeyNoteT{12*(oct-1) + 1, fmt.Sprintf("%dC", oct)}
		m[fmt.Sprintf("Cs%d", oct)] = KeyNoteT{12*(oct-1) + 2, fmt.Sprintf("%dc", oct)}
		m[fmt.Sprintf("Db%d", oct)] = KeyNoteT{12*(oct-1) + 2, fmt.Sprintf("%dc", oct)}
		m[fmt.Sprintf("D%d", oct)] = KeyNoteT{12*(oct-1) + 3, fmt.Sprintf("%dD", oct)}
		m[fmt.Sprintf("Ds%d", oct)] = KeyNoteT{12*(oct-1) + 4, fmt.Sprintf("%dd", oct)}
		m[fmt.Sprintf("Eb%d", oct)] = KeyNoteT{12*(oct-1) + 4, fmt.Sprintf("%dd", oct)}
		m[fmt.Sprintf("E%d", oct)] = KeyNoteT{12*(oct-1) + 5, fmt.Sprintf("%dE", oct)}
		m[fmt.Sprintf("Fb%d", oct)] = KeyNoteT{12*(oct-1) + 5, fmt.Sprintf("%dE", oct)}
		m[fmt.Sprintf("Es%d", oct)] = KeyNoteT{12*(oct-1) + 6, fmt.Sprintf("%dF", oct)}
		m[fmt.Sprintf("F%d", oct)] = KeyNoteT{12*(oct-1) + 6, fmt.Sprintf("%dF", oct)}
		m[fmt.Sprintf("Fs%d", oct)] = KeyNoteT{12*(oct-1) + 7, fmt.Sprintf("%df", oct)}
		m[fmt.Sprintf("Gb%d", oct)] = KeyNoteT{12*(oct-1) + 7, fmt.Sprintf("%df", oct)}
		m[fmt.Sprintf("G%d", oct)] = KeyNoteT{12*(oct-1) + 8, fmt.Sprintf("%dG", oct)}
		m[fmt.Sprintf("Gs%d", oct)] = KeyNoteT{12*(oct-1) + 9, fmt.Sprintf("%dg", oct)}
		m[fmt.Sprintf("Ab%d", oct)] = KeyNoteT{12*(oct-1) + 9, fmt.Sprintf("%dg", oct)}
		m[fmt.Sprintf("A%d", oct)] = KeyNoteT{12*(oct-1) + 10, fmt.Sprintf("%dA", oct)}
		m[fmt.Sprintf("As%d", oct)] = KeyNoteT{12*(oct-1) + 11, fmt.Sprintf("%da", oct)}
		m[fmt.Sprintf("Bb%d", oct)] = KeyNoteT{12*(oct-1) + 11, fmt.Sprintf("%da", oct)}
		m[fmt.Sprintf("B%d", oct)] = KeyNoteT{12*(oct-1) + 12, fmt.Sprintf("%dB", oct)}
	}

	return m
}
