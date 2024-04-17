package physical

type Encoder interface {
	Encode([]Bit) []Bit
}
