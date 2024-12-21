package main

import (
	"sync"

	"github.com/lmatte7/gomesh/github.com/meshtastic/gomeshproto"
)

// Must be initialized with a max size
type buf_packet struct {
	maxSize int
	mux     sync.Mutex
	buf     []gomeshproto.FromRadio_Packet
}

var receiveBuffer buf_packet
var sendBuffer buf_packet

// init
func init() {
	receiveBuffer.SetMaxSize(2048)
	sendBuffer.SetMaxSize(2048)
}

// Set maximum buffer size
func (b *buf_packet) SetMaxSize(i int) {
	if i >= 0 {
		b.maxSize = i
	}
}

// Return max buffer size
func (b *buf_packet) MaxSize() int { return b.maxSize }

// Return current buffer size
func (b *buf_packet) Size() int { return len(b.buf) }

// Add strings to the buffer
func (b *buf_packet) Push(z gomeshproto.FromRadio_Packet) {
	// Lock buffer
	b.mux.Lock()
	if b.Size() < b.MaxSize() {
		// Append
		b.buf = append(b.buf, z)
	}
	// Unlock buffer
	b.mux.Unlock()
}

// Remove an item from the buffer
func (b *buf_packet) Pop() gomeshproto.FromRadio_Packet {
	// Lock buffer
	b.mux.Lock()
	var z gomeshproto.FromRadio_Packet
	// Return empty string if empty
	if b.Size() > 0 {
		// Remove from buffer
		z, b.buf = b.buf[0], b.buf[1:]
	}
	// Unlock buffer
	b.mux.Unlock()
	return z
}
