package physical

import "time"

type DeviceFactory struct {}

type TimeTickerFactory struct {}

type TimeTicker struct {
	ticker *time.Ticker
}

func (tf *TimeTickerFactory) NewTicker(clockPeriod time.Duration) Ticker {
	return &TimeTicker{
		ticker: time.NewTicker(clockPeriod),
	}
}

func (tt *TimeTicker) GetChannel() <-chan time.Time {
	return tt.ticker.C
}

func (df *DeviceFactory) CreateDevice(bitHandler BitHandler) *Device {
	tf := &TimeTickerFactory{}
	sendClockPeriod := time.Millisecond
	listenClockPeriod := time.Nanosecond * 100_000
	encoder := NewManchesterEncoder()
	decoderFactory := &ManchesterDecoderFactory{}
	return NewDevice(tf, sendClockPeriod, listenClockPeriod, encoder, decoderFactory, bitHandler)
}
