package gowl

import (
	"bytes"
	"fmt"
)

// Message represents an SMTP message.
type Message struct {
	Header *Header
	Root   *Part
}

// Render renders the message into bytes in an SMTP format.
func (m *Message) Render() ([]byte, error) {
	buf := bytes.Buffer{}

	head, err := m.Header.Render()
	if err != nil {
		return nil, fmt.Errorf("failed to render message header: %w", err)
	}

	root, err := m.Root.Render()
	if err != nil {
		return nil, fmt.Errorf("failed to render message root part: %w", err)
	}

	buf.Write(head)
	buf.Write(root)

	return buf.Bytes(), nil
}
