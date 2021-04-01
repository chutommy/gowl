package gowl

import (
	"bytes"
	"fmt"
	"io"
)

// Part is a representation of a single piece of SMTP data block which might
// contain a content or in case of multipart another parts.
type Part struct {
	header  *Header
	content io.Reader
	parts   []*Part
}

// NewPart is a constructor of the Part.
func NewPart(header *Header, content io.Reader, parts []*Part) *Part {
	return &Part{
		header:  header,
		content: content,
		parts:   parts,
	}
}

// Header returns a header of the Part.
func (p *Part) Header() *Header {
	return p.header
}

// Content returns a content of the Part as an io.Reader.
func (p *Part) Content() io.Reader {
	return p.content
}

// Parts returns sub-parts of the Part.
func (p *Part) Parts() []*Part {
	return p.parts
}

// SetHeader replaces a header of the Part with the given Header.
func (p *Part) SetHeader(header *Header) {
	p.header = header
}

// SetContent replaces a content of the Part with the given io.Reader.
func (p *Part) SetContent(content io.Reader) {
	p.content = content
}

// SetParts replaces sub-parts of the Part with the given slice of Parts.
func (p *Part) SetParts(parts []*Part) {
	p.parts = parts
}

// Render renders the content of the Part into bytes. It returns a formatted SMTP message Part.
func (p *Part) Render() ([]byte, error) {
	buf := bytes.Buffer{}

	head, err := p.header.Render()
	if err != nil {
		return nil, fmt.Errorf("failed to render part header: %w", err)
	}

	buf.Write(head)

	if p.content != nil {
		buf.Write([]byte{'\n', '\n'})

		if _, err := buf.ReadFrom(p.content); err != nil {
			return nil, fmt.Errorf("failed to read part content: %w", err)
		}
	}

	if p.parts != nil {
		boundary, err := p.header.Boundary()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve sub-part boundary: %w", err)
		}

		bStart := append([]byte{'-', '-'}, boundary...)
		bEnd := append(bStart, '-', '-')

		for _, p := range p.parts {
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
