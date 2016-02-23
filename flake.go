// Flake generates unique identifiers that are roughly sortable by time. Flake can
// run on a cluster of machines and still generate unique IDs without requiring
// worker coordination.
//
// A Flake ID is a 64-bit integer will the following components:
//  - 41 bits is the timestamp with millisecond precision
//  - 10 bits is the host id (uses IP modulo 2^10)
//  - 13 bits is an auto-incrementing sequence for ID requests within the same millisecond
//
// Note: In order to make a millisecond timestamp fit within 41 bits, a custom
// epoch of Jan 1, 2014 00:00:00 is used.

package flake

import (
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

//-----------------------------------------------------------------------------

const (
	HostBits     = 10
	SequenceBits = 13
)

var (
	// Custom Epoch so the timestamp can fit into 41 bits.
	// Jan 1, 2014 00:00:00 UTC
	Epoch       time.Time = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
	MaxHostId   uint64    = (1 << HostBits) - 1
	MaxSequence uint64    = (1 << SequenceBits) - 1
)

// Id represents a unique k-ordered Id
type Id uint64

// String formats the Id as a 16 character hexadecimal string
func (id *Id) String() string {
	return strconv.FormatUint(uint64(*id), 16)
}

// Uint64 formats the Id as an unsigned integer
func (id *Id) Uint64() uint64 {
	return uint64(*id)
}

// Flake is a unique Id generator
type Flake struct {
	prevTime uint64
	hostId   uint64
	sequence uint64
	mu       sync.Mutex
}

// New returns a new Id generator and a possible error condition
func New() (*Flake, error) {
	hostId, err := getHostId()

	if err != nil {
		return nil, err
	}

	return &Flake{
		sequence: 0,
		prevTime: getTimestamp(),
		hostId:   hostId,
	}, nil
}

// NextId returns a new Id from the generator
func (f *Flake) NextId() Id {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := getTimestamp()

	if now < f.prevTime {
		now = f.prevTime
	}

	// Use the sequence number if the id request is in the same millisecond as
	// the previous request.
	if now == f.prevTime {
		f.sequence += 1
	} else {
		f.sequence = 0
	}

	// Bump the timestamp by 1ms if we run out of sequence bits.
	if f.sequence > MaxSequence {
		now += 1
		f.sequence = 0
	}

	f.prevTime = now

	timestamp := now << (HostBits + SequenceBits)
	hostid := f.hostId << SequenceBits

	return Id(timestamp | hostid | f.sequence)
}

// getTimestamp returns the timestamp in milliseconds adjusted for the custom
// epoch
func getTimestamp() uint64 {
	return uint64(time.Since(Epoch).Nanoseconds() / 1e6)
}

// getHostId returns the host id using the IP address of the machine
func getHostId() (uint64, error) {
	h, err := os.Hostname()

	if err != nil {
		return 0, err
	}

	addrs, err := net.LookupIP(h)

	if err != nil {
		return 0, err
	}

	a := addrs[0].To4()
	if a == nil {
		a = addrs[0].To16()
	}

	ip := uint64((a[0] * 1 << 24) + (a[1] * 1 << 16) + (a[2] * 1 << 8) + a[3])

	return ip % MaxHostId, nil
}
