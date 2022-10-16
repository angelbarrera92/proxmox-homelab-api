package probes

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type PortProber struct {
	Host string
	Port int
}

func (p PortProber) Probe() (bool, error) {
	timeout := 3 * time.Second
	hostAndPort := net.JoinHostPort(p.Host, strconv.Itoa(p.Port))
	conn, err := net.DialTimeout("tcp", hostAndPort, timeout)
	if err != nil {
		return false, fmt.Errorf("error dialing port %d: %w", p.Port, err)
	}
	if conn != nil {
		defer conn.Close()
		return true, nil
	}
	return false, nil
}
