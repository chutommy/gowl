package gowl

import (
	"bytes"
	"fmt"
	"io"
)

// Part is a representation of a single piece of SMTP data block which might
// contain a content or in case of multipart another parts.
type Part struct {
	Header  *Header
	Content io.Reader
	Parts   []*Part
}

// Render renders the content of the Part into bytes. It returns a formatted SMTP message Part.
func (p *Part) Render() ([]byte, error) {
	buf := bytes.Buffer{}

	head, err := p.Header.Render()
	if err != nil {
		return nil, fmt.Errorf("failed to render part header: %w", err)
	}

	buf.Write(head)

	if p.Content != nil {
		buf.Write([]byte{'\n', '\n'})

		if _, err := buf.ReadFrom(p.Content); err != nil {
			return nil, fmt.Errorf("failed to read part content: %w", err)
		}
	}

	if p.Parts != nil {
		boundary, err := p.Header.Boundary()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve sub-part boundary: %w", err)
		}

		bStart := append([]byte{'-', '-'}, boundary...)
		bEnd := append(bStart, '-', '-')

		for _, p := range p.Parts {
			buf.Write([]byte{'\n', '\n'})
			buf.Write(bStart)
			buf.WriteRune('\n')

			part, err := p.Render()
			if err != nil {
				return nil, fmt.Errorf("failed to render sub-part part: %w", err)
			}

			buf.Write(part)
		}

		buf.Write([]byte{'\n', '\n'})
		buf.Write(bEnd)
	}

	return buf.Bytes(), nil
}
