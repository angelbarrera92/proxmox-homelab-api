package probes

type ProxmoxProber struct {
	Host   string
	Port   int
	VMName string
}

func (p ProxmoxProber) Probe() (bool, error) {
	return false, nil
}
