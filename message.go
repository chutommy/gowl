package gowl

import (
	"bytes"
	"fmt"
)

// Message represents an SMTP message.
type Message struct {
	header   *Header
	rootPart *Part
}

// NewMessage is a constructor of the Message.
func NewMessage(header *Header, rootPart *Part) *Message {
	return &Message{
		header:   header,
		rootPart: rootPart,
	}
}

// Header returns the header of the Message.
func (m *Message) Header() *Header {
	return m.header
}

// RootPart returns the part at the root of the Message.
func (m *Message) RootPart() *Part {
	return m.rootPart
}

// SetHeader replaces the header of the Message with the given Header.
func (m *Message) SetHeader(header *Header) {
	m.header = header
}

// SetHeader replaces the part at the root of the Message with the given Part.
func (m *Message) SetRootPart(rootPart *Part) {
	m.rootPart = rootPart
}

// Render renders the message into bytes in an SMTP format.
func (m *Message) Render() ([]byte, error) {
	buf := bytes.Buffer{}

	head, err := m.header.Render()
	if err != nil {
		return nil, fmt.Errorf("failed to render message header: %w", err)
	}

	root, err := m.rootPart.Render()
	if err != nil {
		return nil, fmt.Errorf("failed to render message root part: %w", err)
	}

	buf.Write(head)
	buf.Write(root)

	return buf.Bytes(), nil
}
