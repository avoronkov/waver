package wav

func CreateDefaultWav() *Wav {
	return &Wav{
		Fmt: &WavFmt{
			CompressionCode:          1,
			NumberOfChannels:         2,
			SampleRate:               44100,
			AvgBps:                   44100 * 16 / 8 * 2,
			BlockAlign:               4,
			SignificantBitsPerSample: 16,
		},
		Data: &Data{},
	}
}
