package supervisor

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
	"testing"
)

// TestHelperProcess is the entrypoint for the mock worker process.
// It is invoked by the Supervisor tests via exec.Command when GO_WANT_HELPER_PROCESS=1.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// Prevent the test harness from complaining about no tests run, and exit cleanly
	defer os.Exit(0)

	mode := os.Getenv("POGO_MOCK_WORKER_MODE")

	if mode == "crash_immediate" {
		os.Exit(1)
	}

	// FD 3: Input (From Host)
	// FD 4: Output (To Host)
	// These are mapped by exec.Cmd.ExtraFiles
	in := os.NewFile(3, "pipe_in")
	out := os.NewFile(4, "pipe_out")

	if in == nil || out == nil {
		return
	}

	if mode == "hang_handshake" {
		// Read header but never respond, simulating a hung process
		buf := make([]byte, 5)
		_, _ = io.ReadFull(in, buf)
		select {} // Block forever
	}

	// --- Handshake ---
	if err := handleHandshake(in, out); err != nil {
		return
	}

	// --- Main Loop ---
	for {
		// Read Header
		header := make([]byte, 5)
		_, err := io.ReadFull(in, header)
		if err != nil {
			return // Host closed connection
		}

		length := binary.BigEndian.Uint32(header[0:4])
		pktType := header[4]

		// Read Body
		body := make([]byte, length)
		_, err = io.ReadFull(in, body)
		if err != nil {
			return
		}

		if pktType == PktTypeShutdown {
			return
		}

		if pktType == PktTypeData {
			var task map[string]any
			_ = json.Unmarshal(body, &task)

			// Echo logic
			response := map[string]any{
				"status": "success",
				"result": map[string]any{
					"echo": task["payload"],
					"pid":  os.Getpid(),
				},
			}

			_ = sendJson(out, PktTypeData, response)
		}
	}
}

func handleHandshake(in io.Reader, out io.Writer) error {
	// Read Hello
	header := make([]byte, 5)
	if _, err := io.ReadFull(in, header); err != nil {
		return err
	}
	length := binary.BigEndian.Uint32(header[0:4])
	// pktType := header[4] (Should be Hello)

	body := make([]byte, length)
	if _, err := io.ReadFull(in, body); err != nil {
		return err
	}

	// Send Ack
	ack := map[string]any{
		"type": "HELLO_ACK",
		"capabilities": map[string]any{
			"protocol": "json",
			"shm":      false,
		},
	}
	return sendJson(out, PktTypeHello, ack)
}

func sendJson(w io.Writer, typ byte, data any) error {
	b, _ := json.Marshal(data)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(b)))

	if _, err := w.Write(lenBuf); err != nil {
		return err
	}
	if _, err := w.Write([]byte{typ}); err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
