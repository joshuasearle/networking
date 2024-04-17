package physical

import "sync"

type Cable struct {
	mu  sync.Mutex
	bit Bit

	iface1 Iface
	iface2 Iface

	device1 *Device
	device2 *Device
}

func NewCable() *Cable {
	return &Cable{}
}

func (c *Cable) Read() Bit {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.bit
}

func (c *Cable) Write(b Bit) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bit = b
}
