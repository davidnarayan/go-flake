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
	"strconv"
	"sync"
	"time"
)

//-----------------------------------------------------------------------------

const (
	hostBits     = 10
	sequenceBits = 13
)

var (
	// Custom Epoch so the timestamp can fit into 41 bits.
	// Jan 1, 2014 00:00:00 UTC
	epoch              = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
	maxHostID   uint64 = (1 << hostBits) - 1
	maxSequence uint64 = (1 << sequenceBits) - 1
)

// ID represents a unique k-ordered Id
type ID uint64

// String formats the Id as a 16 character hexadecimal string
func (id *ID) String() string {
	return strconv.FormatUint(uint64(*id), 16)
}

// Uint64 formats the Id as an unsigned integer
func (id *ID) Uint64() uint64 {
	return uint64(*id)
}

// Flake is a unique Id generator
type Flake struct {
	prevTime uint64
	HostID   uint64
	sequence uint64
	mu       sync.Mutex
}

// NextId returns a new ID from the generator
func (f *Flake) NextId() ID {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := getTimestamp()

	if now < f.prevTime {
		now = f.prevTime
	}

	// Use the sequence number if the id request is in the same millisecond as
	// the previous request.
	if now == f.prevTime {
		f.sequence++
	} else {
		f.sequence = 0
	}

	// Bump the timestamp by 1ms if we run out of sequence bits.
	if f.sequence > maxSequence {
		now++
		f.sequence = 0
	}

	f.prevTime = now

	timestamp := now << (hostBits + sequenceBits)
	HostID := f.HostID << sequenceBits

	return ID(timestamp | HostID | f.sequence)
}

// getTimestamp returns the timestamp in milliseconds adjusted for the custom
// epoch
func getTimestamp() uint64 {
	return uint64(time.Since(epoch).Nanoseconds() / 1e6)
}
