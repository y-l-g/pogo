package supervisor

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// ErrPayloadTooLarge indicates a packet exceeded MaxPayloadSize.
var ErrPayloadTooLarge = errors.New("payload too large")

// Transport handles the low-level IPC protocol (5-byte header + body).
type Transport struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

// NewTransport creates a new Transport instance.
func NewTransport(r io.ReadCloser, w io.WriteCloser) *Transport {
	return &Transport{
		reader: r,
		writer: w,
	}
}

// Close closes the underlying reader and writer.
func (t *Transport) Close() error {
	rErr := t.reader.Close()
	wErr := t.writer.Close()
	if rErr != nil {
		return rErr
	}
	return wErr
}

// SetWriteDeadline sets the write deadline for the underlying connection.
func (t *Transport) SetWriteDeadline(deadline time.Time) error {
	if f, ok := t.writer.(*os.File); ok {
		return f.SetWriteDeadline(deadline)
	}
	return nil
}

// SetReadDeadline sets the read deadline for the underlying connection.
func (t *Transport) SetReadDeadline(deadline time.Time) error {
	if f, ok := t.reader.(*os.File); ok {
		return f.SetReadDeadline(deadline)
	}
	return nil
}

// WritePacket sends a packet with the given type and body.
// Header format: [Length (4 bytes)][Type (1 byte)]
func (t *Transport) WritePacket(pktType byte, body []byte) error {
	length := uint32(len(body))
	header := make([]byte, 5)

	binary.BigEndian.PutUint32(header[0:4], length)
	header[4] = pktType

	// Write Header
	if _, err := t.writer.Write(header); err != nil {
		return err
	}

	// Write Body
	if length > 0 {
		if _, err := t.writer.Write(body); err != nil {
			return err
		}
	}

	return nil
}

// ReadPacket reads a packet and returns its type and body.
// It enforces MaxPayloadSize.
func (t *Transport) ReadPacket() (byte, []byte, error) {
	header := make([]byte, 5)
	if _, err := io.ReadFull(t.reader, header); err != nil {
		return 0, nil, err
	}

	length := binary.BigEndian.Uint32(header[0:4])
	pktType := header[4]

	if length > MaxPayloadSize {
		return pktType, nil, fmt.Errorf("%w: %d > %d", ErrPayloadTooLarge, length, MaxPayloadSize)
	}

	if length == 0 {
		return pktType, []byte{}, nil
	}

	body := make([]byte, length)
	if _, err := io.ReadFull(t.reader, body); err != nil {
		return pktType, nil, err
	}

	return pktType, body, nil
}
