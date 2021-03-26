package gowl

import (
	"bytes"
	"errors"
	"strings"
)

// Error codes returned by failures to render to SMTP.
var (
	ErrNoValues   = errors.New("the attribute Values of Field is empty")
	ErrNoBoundary = errors.New("the Header has no Content-Type field with boundary parameter")
)

// Header represents SMTP Header.
type Header struct {
	Fields []*Field
}

// Render renders the Header fields and returns them in bytes.
// It writes each field on its own line.
func (h *Header) Render() ([]byte, error) {
	fs := make([][]byte, len(h.Fields))
	var err error

	// render and set each field
	for i, f := range h.Fields {
		if fs[i], err = f.Render(); err != nil {
			return nil, err
		}
	}

	return bytes.Join(fs, []byte{'\n'}), nil
}

// Boundary queries the Header and tries to find a boundary of the Content-Type.
// If there's no boundary parameter inside Content-Type ErrNoBoundary is returned.
func (h *Header) Boundary() ([]byte, error) {
	// iterate over fields
	for _, f := range h.Fields {
		if f.Name == "Content-Type" {
			// iterate over values
			for _, v := range f.Values {
				if strings.Contains(v, "boundary=") {
					// boundary found
					return []byte(strings.Trim(strings.TrimPrefix(v, "boundary="), "\"")), nil
				}
			}
		}
	}

	return nil, ErrNoBoundary
}

// Field represents a single SMTP Header field.
type Field struct {
	Name   string
	Values []string
}

// Render renders the content of the field into bytes. It returns formatted
// SMTP Header Field. The Field's Values are separated with semicolons.
func (f *Field) Render() ([]byte, error) {
	if len(f.Values) == 0 {
		return nil, ErrNoValues
	}
	return []byte(f.Name + ": " + strings.Join(f.Values, "; ")), nil
}
