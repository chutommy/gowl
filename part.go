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

// Render renders the Part's content into bytes. It returns formatted SMTP message Part.
func (p *Part) Render() ([]byte, error) {
	buf := bytes.Buffer{}

	head, err := p.Header.Render()
	if err != nil {
		return nil, fmt.Errorf("failed to render header: %w", err)
	}
	buf.Write(head)

	if p.Content != nil {
		buf.Write([]byte{'\n', '\n'})

		if _, err := buf.ReadFrom(p.Content); err != nil {
			return nil, fmt.Errorf("failed to read content: %w", err)
		}
	}

	if p.Parts != nil {
		boundary, err := p.Header.Boundary()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve boundary: %w", err)
		}
		bStart := append([]byte{'-', '-'}, boundary...)
		bEnd := append(bStart, '-', '-')

		for _, p := range p.Parts {
			buf.Write([]byte{'\n', '\n'})
			buf.Write(bStart)
			buf.WriteRune('\n')

			part, err := p.Render()
			if err != nil {
				return nil, fmt.Errorf("failed to render a part: %w", err)
			}
			buf.Write(part)
		}

		buf.Write([]byte{'\n', '\n'})
		buf.Write(bEnd)
	}

	return buf.Bytes(), nil
}
