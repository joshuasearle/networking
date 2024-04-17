package physical

type Bit bool

const (
	Zero Bit = false
	One  Bit = true
)

func (b Bit) String() string {
	if b {
		return "1"
	}
	return "0"
}

func (b Bit) Opposite() Bit {
	if b == Zero {
		return One
	}
	return Zero
}
