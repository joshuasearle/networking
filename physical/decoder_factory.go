package physical

type IfaceBitHandler interface {
	Handle(Bit)
}

type DecoderFactory interface {
	NewDecoder(IfaceBitHandler) Decoder
}
