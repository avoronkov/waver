package multiplayer

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"sync"
	"time"

	"github.com/avoronkov/waver/lib/midisynth/wav"
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type MultiPlayer struct {
	settings *wav.Settings
	dt       float64

	playingWaves []playingWave

	startTime   time.Time
	sampleCount int64

	mutex sync.Mutex
}

func New(settings *wav.Settings) *MultiPlayer {
	return &MultiPlayer{
		settings:  settings,
		dt:        1.0 / float64(settings.SampleRate),
		startTime: time.Now(),
	}
}

type playingWave struct {
	wave    waves.Wave
	noteCtx *waves.NoteCtx
	time    float64
}

const maxInt16Amp = (1 << 15) - 1

func (m *MultiPlayer) Read(data []byte) (n int, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	l := 640
	if len(data) < l {
		l = len(data)
	}

	buff := new(bytes.Buffer)
	for buff.Len() < l {
		finishedWaves := []int{}
		value := 0.0
		for i, pw := range m.playingWaves {
			if pw.time >= 0.0 {
				waveValue := pw.wave.Value(pw.time, pw.noteCtx)
				if math.IsNaN(waveValue) {
					finishedWaves = append(finishedWaves, i)
				} else {
					value += waveValue
				}
			}
			m.playingWaves[i].time += m.dt
		}
		intValue := int16(maxInt16Amp * value)
		for ch := 0; ch < m.settings.ChannelNum; ch++ {
			_ = binary.Write(buff, binary.LittleEndian, intValue)
		}
		m.sampleCount++
		m.removePlayingWaves(finishedWaves)
	}
	n = copy(data, buff.Bytes())
	return n, nil
}

func (m *MultiPlayer) AddWaveAt(at time.Time, wave waves.Wave, noteCtx *waves.NoteCtx) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delay := m.countDelay(at)
	tm := -delay
	if tm > 0 {
		log.Printf("[WARN] slow processing, delay = %v", tm)
		tm = 0
	}
	m.playingWaves = append(m.playingWaves, playingWave{
		wave:    wave,
		noteCtx: noteCtx,
		time:    tm,
	})
}

func (m *MultiPlayer) countDelay(to time.Time) float64 {
	seconds := float64(m.sampleCount) / float64(m.settings.SampleRate)
	currentTime := m.startTime.Add(time.Duration(seconds * float64(time.Second)))
	dur := to.Sub(currentTime)
	return float64(dur) / float64(time.Second)
}

func (m *MultiPlayer) removePlayingWaves(indexes []int) {
	l := len(indexes)
	last := len(m.playingWaves) - 1
	for i := l - 1; i >= 0; i-- {
		// Remove the element at index i from a.
		idx := indexes[i]
		m.playingWaves[idx] = m.playingWaves[last] // Copy last element to index i.
		m.playingWaves[last] = playingWave{}       // Erase last element (write zero value).
		m.playingWaves = m.playingWaves[:last]     // Truncate slice.
		last--
	}
}
