package ipc

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"log/slog"
)

// Client handles communication with the LabVIEW bridge process.
// It sends requests over an io.Writer and receives responses over an io.Reader.
type Client struct {
	writer io.Writer
	reader io.Reader
	mu     sync.Mutex

	// Response handling
	pending   map[string]chan *Response
	pendingMu sync.RWMutex

	// Request ID generation
	nextID atomic.Uint64

	// Background reader
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

// NewClient creates a new IPC client that communicates via the provided reader/writer.
// For stdio transport, use os.Stdin/os.Stdout.
// For TCP transport, use net.Conn.
func NewClient(ctx context.Context, reader io.Reader, writer io.Writer) *Client {
	clientCtx, cancel := context.WithCancel(ctx)

	c := &Client{
		writer:  writer,
		reader:  reader,
		pending: make(map[string]chan *Response),
		ctx:     clientCtx,
		cancel:  cancel,
		done:    make(chan struct{}),
	}

	// Start background goroutine to read responses
	go c.readResponses()

	return c
}

// Call sends a request to the LabVIEW bridge and waits for the response.
// It returns the response payload or an error if the operation failed.
func (c *Client) Call(ctx context.Context, tool string, input interface{}) (json.RawMessage, error) {
	// Generate unique request ID
	id := fmt.Sprintf("%d", c.nextID.Add(1))

	// Create request
	req, err := NewRequest(id, tool, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Create response channel
	respChan := make(chan *Response, 1)
	c.pendingMu.Lock()
	c.pending[id] = respChan
	c.pendingMu.Unlock()

	// Ensure cleanup
	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
		close(respChan)
	}()

	// Send request
	if err := c.sendRequest(req); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Wait for response
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.ctx.Done():
		return nil, fmt.Errorf("client closed")
	case resp := <-respChan:
		if resp.Error != "" {
			return nil, fmt.Errorf("bridge error: %s", resp.Error)
		}
		return resp.Payload, nil
	}
}

// sendRequest sends a request to the bridge as line-delimited JSON.
func (c *Client) sendRequest(req *Request) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Write line-delimited JSON
	if _, err := c.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}
	if _, err := c.writer.Write([]byte("\n")); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	slog.Debug("IPC request sent", slog.String("id", req.ID), slog.String("tool", req.Tool))
	return nil
}

// readResponses is a background goroutine that reads responses from the bridge.
func (c *Client) readResponses() {
	defer close(c.done)

	scanner := bufio.NewScanner(c.reader)
	for scanner.Scan() {
		line := scanner.Bytes()

		var resp Response
		if err := json.Unmarshal(line, &resp); err != nil {
			slog.Error("Failed to unmarshal response", slog.Any("error", err))
			continue
		}

		// Route response to waiting caller
		c.pendingMu.RLock()
		respChan, ok := c.pending[resp.ID]
		c.pendingMu.RUnlock()

		if ok {
			select {
			case respChan <- &resp:
				slog.Debug("IPC response received", slog.String("id", resp.ID))
			case <-c.ctx.Done():
				return
			}
		} else {
			slog.Warn("Received response for unknown request ID", slog.String("id", resp.ID))
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Scanner error", slog.Any("error", err))
	}
}

// Close stops the client and releases resources.
func (c *Client) Close() error {
	c.cancel()
	<-c.done
	return nil
}
