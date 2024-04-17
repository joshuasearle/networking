package physical

type Decoder interface {
	Handle(Bit)
}
