package wav

type Settings struct {
	SampleRate      int
	ChannelNum      int
	BitDepthInBytes int
}

var Default = Settings{
	SampleRate:      44000,
	ChannelNum:      2,
	BitDepthInBytes: 2,
}
