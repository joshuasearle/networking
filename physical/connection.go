package physical

type CableFull struct{}

func (c *CableFull) Error() string {
	return "Cable is full"
}

type DeviceNotConnected struct{}

func (d *DeviceNotConnected) Error() string {
	return "Device not connected"
}

func Connect(c *Cable, d *Device, i Iface) error {
	err := d.connect(c, i)
	if err != nil {
		return err
	}

	if c.device1 == nil {
		c.device1 = d
		c.iface1 = i

		return nil
	}

	if c.device2 == nil {

		c.device2 = d
		c.iface2 = i
		return nil
	}

	return &CableFull{}
}

func Disconnect(c *Cable, d *Device) error {
	err := d.disconnect(c)
	if err != nil {
		return err
	}

	if c.device1 == d {
		c.device1 = nil
		c.iface1 = NewIface("")
		return nil
	}

	if c.device2 == d {
		c.device2 = nil
		c.iface2 = NewIface("")
		return nil
	}

	return &DeviceNotConnected{}
}
