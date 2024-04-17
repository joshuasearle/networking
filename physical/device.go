package physical

import "time"

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
	ifacePreviousBit                    map[Iface]Bit
	ifaceClockPeriodsSinceLastBitChange map[Iface]int
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
		ifacePreviousBit:                    make(map[Iface]Bit),
		ifaceClockPeriodsSinceLastBitChange: make(map[Iface]int),
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
	d.ifacePreviousBit[i] = Zero
	d.ifaceClockPeriodsSinceLastBitChange[i] = 0

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
			delete(d.ifacePreviousBit, i)
			delete(d.ifaceClockPeriodsSinceLastBitChange, i)
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
		prev := d.ifacePreviousBit[i]
		d.ifacePreviousBit[i] = curBit

		// Set the clock period count
		if curBit == prev {
			d.ifaceClockPeriodsSinceLastBitChange[i]++
		} else {
			d.ifaceClockPeriodsSinceLastBitChange[i] = 0
		}

		// Handle bit change
		if curBit != prev {
			decoder := d.ifaceDecoders[i]
			decoder.Decode(curBit)
		}
	}
}
