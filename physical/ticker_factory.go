package physical

import "time"

type TickerFactory interface {
	NewTicker(time.Duration) Ticker
}
