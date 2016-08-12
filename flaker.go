package flake

import (
	"errors"
	"net"
	"os"
)

// Flaker interface to make it abstract
type Flaker interface {
	NextId() ID
}

// New returns a new Flake instance and a possible error condition
func New() (*Flake, error) {
	HostID, err := getHostID()

	if err != nil {
		return nil, err
	}

	return &Flake{
		sequence: 0,
		prevTime: getTimestamp(),
		HostID:   HostID,
	}, nil
}

// getHostID returns the host id using the IP address of the machine
func getHostID() (uint64, error) {
	h, err := os.Hostname()

	if err != nil {
		return 0, err
	}

	addrs, err := net.LookupIP(h)
	a := addrs[0]
	startPos := len(a) - 4
	if startPos < 0 {
		return 0, errors.New("invalid local IP address " + a.String())
	}
	ip := (uint64(a[startPos]) << 24) + (uint64(a[startPos+1]) << 16) + (uint64(a[startPos+2]) << 8) + uint64(a[startPos+3])

	return ip % maxHostID, nil
}
