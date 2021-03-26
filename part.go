package gowl

import (
	"bytes"
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

	// write header
	head, err := p.Header.Render()
	if err != nil {
		return nil, err
	}
	buf.Write(head)

	// write content
	if cnt := p.Content; cnt != nil {
		buf.Write([]byte{'\n', '\n'})
		if _, err := buf.ReadFrom(p.Content); err != nil {
			return nil, err
		}
	}

	if p.Parts != nil {
		// get boundaries
		bound, err := p.Header.Boundary()
		if err != nil {
			return nil, err
		}
		bA := append([]byte{'-', '-'}, bound...) // beginning boundary
		bB := append(bA, '-', '-')               // ending boundary

		// write included parts
		for _, p := range p.Parts {
			buf.Write([]byte{'\n', '\n'})
			buf.Write(bA)
			buf.WriteRune('\n')

			part, err := p.Render()
			if err != nil {
				return nil, err
			}
			buf.Write(part)
		}

		// write ending boundary
		buf.Write([]byte{'\n', '\n'})
		buf.Write(bB)
	}

	return buf.Bytes(), nil
}
