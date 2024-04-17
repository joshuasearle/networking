package physical

import (
	"time"
)

type IfaceBitForwarder struct {
	bitHandler BitHandler
	iface      Iface
}

func NewIfaceBitForwarder(bitHandler BitHandler, iface Iface) *IfaceBitForwarder {
	return &IfaceBitForwarder{
		bitHandler: bitHandler,
		iface:      iface,
	}
}

func (i *IfaceBitForwarder) Handle(bit Bit) {
	i.bitHandler.Handle(bit, i.iface)
}

type Device struct {
	ifaceCables                         map[Iface]*Cable
	ifaceDecoders                       map[Iface]Decoder

	tickerFactory TickerFactory

	sendClockPeriod   time.Duration
	listenClockPeriod time.Duration

	encoder        Encoder
	decoderFactory DecoderFactory

	bitHandler BitHandler
}

func NewDevice(tickerFactory TickerFactory, sendClockPeriod time.Duration, listenClockPeriod time.Duration, encoder Encoder, decoderFactory DecoderFactory, bitHandler BitHandler) *Device {
	return &Device{
		ifaceCables:                         make(map[Iface]*Cable),
		ifaceDecoders:                       make(map[Iface]Decoder),

		tickerFactory: tickerFactory,

		sendClockPeriod:   sendClockPeriod,
		listenClockPeriod: listenClockPeriod,

		encoder:        encoder,
		decoderFactory: decoderFactory,

		bitHandler: bitHandler,
	}
}

type IfaceInUse struct{}

func (i *IfaceInUse) Error() string {
	return "Iface in use"
}

func (d *Device) connect(c *Cable, i Iface) error {
	if _, ok := d.ifaceCables[i]; ok {
		return &IfaceInUse{}
	}

	d.ifaceCables[i] = c

	ifaceBitForwarder := NewIfaceBitForwarder(d.bitHandler, i)

	decoder := d.decoderFactory.NewDecoder(ifaceBitForwarder)

	d.ifaceDecoders[i] = decoder

	return nil
}

type IfaceNotInUse struct{}

func (i *IfaceNotInUse) Error() string {
	return "Iface not in use"
}

func (d *Device) disconnect(c *Cable) error {
	for i, cable := range d.ifaceCables {
		if cable == c {
			delete(d.ifaceCables, i)
			return nil
		}
	}

	return &IfaceNotInUse{}
}

func (d *Device) Send(i Iface, bits []Bit) {
	t := d.tickerFactory.NewTicker(d.sendClockPeriod)

	encoding := d.encoder.Encode(bits)

	for _, bit := range encoding {
		<-t.GetChannel()
		cable := d.ifaceCables[i]
		cable.Write(bit)
	}
}

func (d *Device) Listen() {
	t := d.tickerFactory.NewTicker(d.listenClockPeriod)

	for {
		<-t.GetChannel()
		d.handleListenClockCycle()
	}
}

func (d *Device) handleListenClockCycle() {
	for i, cable := range d.ifaceCables {
		curBit := cable.Read()

		decoder := d.ifaceDecoders[i]
		decoder.Handle(curBit)
	}
}
