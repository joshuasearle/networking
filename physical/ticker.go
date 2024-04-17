package physical

import "time"

type Ticker interface {
	GetChannel() <-chan time.Time
}
