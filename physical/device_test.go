package physical_test

import (
	"networking/physical"
	"testing"
	"time"
)

type MockTicker struct {
	channel chan time.Time
}

func NewMockTicker() *MockTicker {
	return &MockTicker{
		channel: make(chan time.Time),
	}
}

func (t *MockTicker) GetChannel() <-chan time.Time {
	return t.channel
}

func (t *MockTicker) Tick() {
	t.channel <- time.Now()

	// Sleep to allow the tick to be processed
	time.Sleep(1 * time.Millisecond)
}

type MockTickerFactory struct {
	tickers []physical.Ticker
}

func NewMockTickerFactory(tickers []physical.Ticker) *MockTickerFactory {
	return &MockTickerFactory{
		tickers: tickers,
	}
}

func (tf *MockTickerFactory) NewTicker(time.Duration) physical.Ticker {
	first := tf.tickers[0]
	tf.tickers = tf.tickers[1:]

	return first
}

type MockEncoder struct{}

func NewMockEncoder() *MockEncoder {
	return &MockEncoder{}
}
func (e *MockEncoder) Encode(data []physical.Bit) []physical.Bit {
	return data
}

type MockDecoder struct {
	ifaceBitHandler physical.IfaceBitHandler
}

func NewMockDecoder(ifaceBitHandler physical.IfaceBitHandler) *MockDecoder {
	return &MockDecoder{
		ifaceBitHandler: ifaceBitHandler,
	}
}

func (d *MockDecoder) Handle(bit physical.Bit) {
	d.ifaceBitHandler.Handle(bit)
}

type MockDecoderFactory struct{}

func NewMockDecoderFactory() *MockDecoderFactory {
	return &MockDecoderFactory{}
}

func (df *MockDecoderFactory) NewDecoder(ifaceBitHandler physical.IfaceBitHandler) physical.Decoder {
	decoder := NewMockDecoder(ifaceBitHandler)

	return decoder
}

type MockBitHandler struct {
	ifaceBits map[physical.Iface][]physical.Bit
}

func NewMockBitHandler() *MockBitHandler {
	return &MockBitHandler{
		ifaceBits: make(map[physical.Iface][]physical.Bit),
	}
}

func (bh *MockBitHandler) Handle(bit physical.Bit, i physical.Iface) {
	if _, ok := bh.ifaceBits[i]; !ok {
		bh.ifaceBits[i] = []physical.Bit{}
	}

	bh.ifaceBits[i] = append(bh.ifaceBits[i], bit)
}

func SetupDevice(tickers []physical.Ticker) (*physical.Device, *MockBitHandler) {
	tf := NewMockTickerFactory(tickers)

	e := NewMockEncoder()

	df := NewMockDecoderFactory()

	bh := NewMockBitHandler()

	return physical.NewDevice(tf, 1*time.Second, 100*time.Millisecond, e, df, bh), bh
}

func TestThatDevicesCanConnect(t *testing.T) {
	c := physical.NewCable()

	sendTicker := NewMockTicker()
	d1, _ := SetupDevice([]physical.Ticker{sendTicker})
	i1 := physical.NewIface("eth0")

	listenTicker := NewMockTicker()
	d2, _ := SetupDevice([]physical.Ticker{listenTicker})
	i2 := physical.NewIface("eth1")

	err := physical.Connect(c, d1, i1)
	if err != nil {
		t.Fatalf("Expected nil, got %v", err)
	}
	err = physical.Connect(c, d2, i2)
	if err != nil {
		t.Fatalf("Expected nil, got %v", err)
	}

	err = physical.Disconnect(c, d2)
	if err != nil {
		t.Fatalf("Expected nil, got %v", err)
	}
	err = physical.Disconnect(c, d1)
	if err != nil {
		t.Fatalf("Expected nil, got %v", err)
	}
}

func TestThatDevicesReceivesBits(t *testing.T) {
	sendTicker := NewMockTicker()
	d1, _ := SetupDevice([]physical.Ticker{sendTicker})
	i1 := physical.NewIface("eth0")

	listenTicker := NewMockTicker()
	d2, bh2 := SetupDevice([]physical.Ticker{listenTicker})
	i2 := physical.NewIface("eth1")

	c := physical.NewCable()

	physical.Connect(c, d1, i1)
	physical.Connect(c, d2, i2)

	data := []physical.Bit{physical.One, physical.Zero, physical.One, physical.Zero}

	go d2.Listen()
	go d1.Send(i1, data)

	sendTicker.Tick()
	listenTicker.Tick()
	sendTicker.Tick()
	listenTicker.Tick()
	sendTicker.Tick()
	listenTicker.Tick()
	sendTicker.Tick()
	listenTicker.Tick()

	bits := bh2.ifaceBits[i2]
	if len(bits) != len(data) {
		t.Fatalf("Expected %v, got %v", data, bits)
	}
}

func TestCableInterference(t *testing.T) {
	sendTicker1 := NewMockTicker()
	d1, _ := SetupDevice([]physical.Ticker{sendTicker1})
	i1 := physical.NewIface("eth0")

	sendTicker2 := NewMockTicker()
	d2, _ := SetupDevice([]physical.Ticker{sendTicker2})
	i2 := physical.NewIface("eth1")

	c := physical.NewCable()

	physical.Connect(c, d1, i1)
	physical.Connect(c, d2, i2)

	data1 := []physical.Bit{physical.One, physical.One, physical.One, physical.One}
	data2 := []physical.Bit{physical.Zero, physical.Zero, physical.Zero, physical.Zero}

	go d1.Send(i1, data1)
	go d2.Send(i2, data2)

	if c.Read() != physical.Zero {
		t.Fatalf("Expected %v, got %v", physical.Zero, c.Read())
	}

	sendTicker1.Tick()

	if c.Read() != physical.One {
		t.Fatalf("Expected %v, got %v", physical.One, c.Read())
	}

	sendTicker2.Tick()

	if c.Read() != physical.Zero {
		t.Fatalf("Expected %v, got %v", physical.One, c.Read())
	}
}
