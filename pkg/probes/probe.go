package probes

type Prober interface {
	Probe() (bool, error)
}
