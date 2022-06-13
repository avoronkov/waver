package signals

import "time"

type StartTime struct {
	Start time.Time
}

func (*StartTime) SignalType() string {
	return "start-time"
}
