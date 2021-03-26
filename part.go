package gowl

import (
	"bytes"
	"io"
	"strings"
)

// Part is a representation of a single piece of SMTP data block which might
// contain a content or in case of multipart another parts.
type Part struct {
	Header  Header
	Content io.Reader
	Parts   []*Part
}

// NewPart creates a new instance of Part with given values.
// func NewPart(header *Header, content io.Reader, parts []*Part) *Part {
// 	return &Part{
// 		Header:  header,
// 		Content: content,
// 		Parts:   parts,
// 	}
// }

// Render renders the Part's content into bytes. It returns formatted SMTP message Part.
func (p *Part) Render() []byte {
	buf := bytes.Buffer{}

	// write header
	buf.Write(p.Header.Render())

	// write content
	if c := p.Content; c != nil {
		buf.Write([]byte{'\n', '\n'})
		if _, err := buf.ReadFrom(p.Content); err != nil {
			panic(err)
		}
	}

	if p.Parts != nil {
		// get boundaries
		bA, bB := boundaries(p.Header)

		// write included parts
		for _, p := range p.Parts {
			buf.Write([]byte{'\n', '\n'})
			buf.Write(bA)
			buf.WriteRune('\n')
			buf.Write(p.Render())
		}

		// write ending boundary
		buf.Write([]byte{'\n', '\n'})
		buf.Write(bB)
	}

	return buf.Bytes()
}

// boundaries returns a beginning and ending boundaries according to the Header.
func boundaries(h Header) ([]byte, []byte) {
	b := getBoundary(h)
	bA := append([]byte{'-', '-'}, b...) // beginning boundary
	bB := append(bA, '-', '-')           // ending boundary
	return bA, bB
}

// getBoundary queries the Header and tries to find a boundary of the Content-Type.
// If there's no boundary parameter inside Content-Type the function panics.
func getBoundary(h Header) []byte {
	for _, f := range h.Fields {
		if f.Name == "Content-Type" {
			for _, v := range f.Values {
				if strings.Contains(v, "boundary=") {
					return []byte(strings.Trim(strings.TrimPrefix(v, "boundary="), "\""))
				}
			}
		}
	}

	panic("Content-Type with boundary value is not set")
}
