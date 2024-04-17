package physical

type ManchesterEncoder struct{}

func NewManchesterEncoder() *ManchesterEncoder {
	return &ManchesterEncoder{}
}

func (e *ManchesterEncoder) Encode(bits []Bit) []Bit {
	encoded := make([]Bit, 2*len(bits))

	for i := 0; i < len(bits); i++ {
		bit := bits[i]

		// Bit is represented by a transition
		// The ending value of the transition is the value of the bit
		encoded[2*i] = bit.Opposite()
		encoded[2*i+1] = bit
	}

	return encoded
}
