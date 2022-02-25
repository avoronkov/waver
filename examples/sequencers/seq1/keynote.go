package main

import "fmt"

type KeyNote int64

const (
	C1  KeyNote = 1
	Cs1         = 2
	Db1         = 2
	D1          = 3
	Ds1         = 4
	Eb1         = 4
	E1          = 5
	F1          = 6
	Fs1         = 7
	Gb1         = 7
	G1          = 8
	Gs1         = 9
	Ab1         = 9
	A1          = 10
	As1         = 11
	Bb1         = 11
	B1          = 12

	C2  KeyNote = 13
	Cs2         = 14
	Db2         = 14
	D2          = 15
	Ds2         = 16
	Eb2         = 16
	E2          = 17
	F2          = 18
	Fs2         = 19
	Gb2         = 19
	G2          = 20
	Gs2         = 21
	Ab2         = 21
	A2          = 22
	As2         = 23
	Bb2         = 23
	B2          = 24

	C3  KeyNote = 1 + 24
	Cs3         = 2 + 24
	Db3         = 2 + 24
	D3          = 3 + 24
	Ds3         = 4 + 24
	Eb3         = 4 + 24
	E3          = 5 + 24
	F3          = 6 + 24
	Fs3         = 7 + 24
	Gb3         = 7 + 24
	G3          = 8 + 24
	Gs3         = 9 + 24
	Ab3         = 9 + 24
	A3          = 10 + 24
	As3         = 11 + 24
	Bb3         = 11 + 24
	B3          = 12 + 24

	C4  KeyNote = 1 + 36
	Cs4         = 2 + 36
	Db4         = 2 + 36
	D4          = 3 + 36
	Ds4         = 4 + 36
	Eb4         = 4 + 36
	E4          = 5 + 36
	F4          = 6 + 36
	Fs4         = 7 + 36
	Gb4         = 7 + 36
	G4          = 8 + 36
	Gs4         = 9 + 36
	Ab4         = 9 + 36
	A4          = 10 + 36
	As4         = 11 + 36
	Bb4         = 11 + 36
	B4          = 12 + 36

	C5  KeyNote = 1 + 48
	Cs5         = 2 + 48
	Db5         = 2 + 48
	D5          = 3 + 48
	Ds5         = 4 + 48
	Eb5         = 4 + 48
	E5          = 5 + 48
	F5          = 6 + 48
	Fs5         = 7 + 48
	Gb5         = 7 + 48
	G5          = 8 + 48
	Gs5         = 9 + 48
	Ab5         = 9 + 48
	A5          = 10 + 48
	As5         = 11 + 48
	Bb5         = 11 + 48
	B5          = 12 + 48

	C6  KeyNote = 1 + 60
	Cs6         = 2 + 60
	Db6         = 2 + 60
	D6          = 3 + 60
	Ds6         = 4 + 60
	Eb6         = 4 + 60
	E6          = 5 + 60
	F6          = 6 + 60
	Fs6         = 7 + 60
	Gb6         = 7 + 60
	G6          = 8 + 60
	Gs6         = 9 + 60
	Ab6         = 9 + 60
	A6          = 10 + 60
	As6         = 11 + 60
	Bb6         = 11 + 60
	B6          = 12 + 60
)

func (k KeyNote) String() string {
	switch k {
	case C1:
		return "1C"
	case Cs1:
		return "1c"
	case D1:
		return "1D"
	case Ds1:
		return "1d"
	case E1:
		return "1E"
	case F1:
		return "1F"
	case Fs1:
		return "1f"
	case G1:
		return "1G"
	case Gs1:
		return "1g"
	case A1:
		return "1A"
	case As1:
		return "1a"
	case B1:
		return "1B"

	case C2:
		return "2C"
	case Cs2:
		return "2c"
	case D2:
		return "2D"
	case Ds2:
		return "2d"
	case E2:
		return "2E"
	case F2:
		return "2F"
	case Fs2:
		return "2f"
	case G2:
		return "2G"
	case Gs2:
		return "2g"
	case A2:
		return "2A"
	case As2:
		return "2a"
	case B2:
		return "2B"

	case C3:
		return "3C"
	case Cs3:
		return "3c"
	case D3:
		return "3D"
	case Ds3:
		return "3d"
	case E3:
		return "3E"
	case F3:
		return "3F"
	case Fs3:
		return "3f"
	case G3:
		return "3G"
	case Gs3:
		return "3g"
	case A3:
		return "3A"
	case As3:
		return "3a"
	case B3:
		return "3B"

	case C4:
		return "4C"
	case Cs4:
		return "4c"
	case D4:
		return "4D"
	case Ds4:
		return "4d"
	case E4:
		return "4E"
	case F4:
		return "4F"
	case Fs4:
		return "4f"
	case G4:
		return "4G"
	case Gs4:
		return "4g"
	case A4:
		return "4A"
	case As4:
		return "4a"
	case B4:
		return "4B"

	case C5:
		return "5C"
	case Cs5:
		return "5c"
	case D5:
		return "5D"
	case Ds5:
		return "5d"
	case E5:
		return "5E"
	case F5:
		return "5F"
	case Fs5:
		return "5f"
	case G5:
		return "5G"
	case Gs5:
		return "5g"
	case A5:
		return "5A"
	case As5:
		return "5a"
	case B5:
		return "5B"

	case C6:
		return "6C"
	case Cs6:
		return "6c"
	case D6:
		return "6D"
	case Ds6:
		return "6d"
	case E6:
		return "6E"
	case F6:
		return "6F"
	case Fs6:
		return "6f"
	case G6:
		return "6G"
	case Gs6:
		return "6g"
	case A6:
		return "6A"
	case As6:
		return "6a"
	case B6:
		return "6B"
	default:
		panic(fmt.Errorf("Unknown note: %d", k))
	}
}
