package physical

type Iface string

func NewIface(name string) Iface {
	return Iface(name)
}
