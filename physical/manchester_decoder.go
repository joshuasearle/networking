package physical

type ManchesterDecoderFactory struct{}

func (mdf *ManchesterDecoderFactory) NewDecoder(bitHandler IfaceBitHandler) Decoder {
	return &ManchesterDecoder{
		bitHandler: bitHandler,
	}
}

type ManchesterDecoder struct {
	bitHandler IfaceBitHandler
	inPreable  bool

	// Best guess as to number of ticks between smaller transition
	// When 1010 encoded data, T is time between 1 and 0
	// When consecutive same bits, ticks between is T
	// When alternating bits, ticks between is 2T (preamble bits)
	T float64

	prevBit           Bit
	isFirstBit        bool
	isFirstTransition bool

	ticksSinceLastTransition int

	preambleBitCount       int
	minimumPreableBitCount int

	previousTransitionReset bool
}

func NewManchesterDecoder(bitHandler IfaceBitHandler) *ManchesterDecoder {
	return &ManchesterDecoder{
		bitHandler: bitHandler,
		inPreable:  true,
		// Placeholder until we have received first bit
		prevBit:                  Zero,
		isFirstBit:               true,
		isFirstTransition:        true,
		ticksSinceLastTransition: 0,
		preambleBitCount:         0,
		minimumPreableBitCount:   8,
		previousTransitionReset:  false,
	}
}

func (d *ManchesterDecoder) ResetToPreamble() {
	d.inPreable = true
	d.T = float64(d.ticksSinceLastTransition) / 2
	d.preambleBitCount = 0
	d.previousTransitionReset = false
}

func (d *ManchesterDecoder) MoveToDataState(bit Bit) {
	d.inPreable = false
	d.preambleBitCount = 0
	d.previousTransitionReset = true
	d.prevBit = bit
}

func (d *ManchesterDecoder) SendPreambleBit(bit Bit) {
	d.bitHandler.Handle(bit)
	d.prevBit = bit
	d.preambleBitCount++
}

func (d *ManchesterDecoder) SendDataBit(bit Bit) {
	d.bitHandler.Handle(bit)
	d.prevBit = bit
	d.previousTransitionReset = false
}

func (d *ManchesterDecoder) Handle(bit Bit) {
	if d.isFirstBit {
		d.isFirstBit = false
		d.prevBit = bit
	}

	// We only care about transitions between bits
	if d.prevBit == bit {
		d.ticksSinceLastTransition++
		return
	}

	singleTLower := d.T * 0.75
	singleTUpper := d.T * 1.5

	doubleTLower := singleTUpper
	doubleTUpper := d.T * 2.5

	singleT := singleTLower <= float64(d.ticksSinceLastTransition) && float64(d.ticksSinceLastTransition) <= singleTUpper
	doubleT := doubleTLower < float64(d.ticksSinceLastTransition) && float64(d.ticksSinceLastTransition) <= doubleTUpper

	if d.inPreable {
		if d.T == 0 {
			if d.isFirstTransition {
				// If is first transition, time between transitions is unknown
				d.isFirstTransition = false
			} else {
				// If we have had some ticks, we can make a guess
				// As we are in the preamble, we can assume that the time between transitions is 2T
				// This means our guess should be half of the time between transitions
				d.T = float64(d.ticksSinceLastTransition) / 2
			}

			d.SendPreambleBit(bit)
		} else if singleT {
			// This means we are either receiving our first non-consecutive bit,
			// which means we are out of the preamble or
			// we are receiving a preamble bit with a different timing

			// We can use the heuristic of number of preamble bits so far to determine our confidence in T,
			// and therefore move to the next state if we are confident enough
			if d.preambleBitCount >= d.minimumPreableBitCount {
				// We are confident enough in our guess for T,
				// so we can move to the next state
				d.MoveToDataState(bit)

				// We don't need to send the bit we just received,
				// as transition was a reset transition
			} else {
				// We are not confident enough in our guess for T,
				// so revise current T
				d.ResetToPreamble()

				// Send the bit we just received
				d.SendPreambleBit(bit)
			}
		} else if doubleT {
			// This means we are receiving a preamble bit with the same timing
			d.SendPreambleBit(bit)
		} else {
			// We didn't receive a valid timing for the preamble,
			// so need to revise
			d.ResetToPreamble()

			d.SendPreambleBit(bit)
		}

	} else {
		// We are out of the preamble
		if singleT {
			if d.previousTransitionReset {
				// We are receiving a data bit
				d.SendDataBit(bit)
			} else {
				// We are receiving a resetting bit
				d.previousTransitionReset = true
				d.prevBit = bit
			}
		} else if doubleT {
			if d.previousTransitionReset {
				d.ResetToPreamble()
				d.SendPreambleBit(bit)
			} else {
				d.SendDataBit(bit)
			}
		} else {
			// We didn't receive a valid timing for the data,
			// so need to revise
			d.ResetToPreamble()
			d.SendPreambleBit(bit)
		}
	}

	d.ticksSinceLastTransition = 1
}
