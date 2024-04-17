package physical

type BitHandler interface {
	Handle(Bit, Iface)
}
