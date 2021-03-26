package gowl

import (
	"bytes"
	"strings"
)

// Header represents SMTP Header.
type Header struct {
	Fields []Field
}

// NewHeader builds a new empty (non-nil) Header.
// func NewHeader() *Header {
// 	return new(Header)
// }

// Add appends a Field to the end of the Header.
// func (h Header) Add(field Field) {
// 	if h.Fields == nil {
// 		h.Fields = []Field{}
// 	}
// 	h.Fields = append(h.Fields, field)
// }

// Render renders the Header fields and returns them in bytes.
// It writes each field on its own line.
func (h Header) Render() []byte {
	bb := make([][]byte, len(h.Fields))
	for i, f := range h.Fields {
		bb[i] = f.Render()
	}

	return bytes.Join(bb, []byte{'\n'})
}

// Field represents a single SMTP Header field.
type Field struct {
	Name   string
	Values []string
}

// NewField constructs a new Field and set the Name to name and
// its Values to values.
// func NewField(name string, values ...string) Field {
// 	return Field{
// 		Name:   name,
// 		Values: values,
// 	}
// }

// Render renders the content of the field into bytes. It returns formatted
// SMTP Header Field. The Field's Values are separated with semicolons.
func (f Field) Render() []byte {
	return []byte(f.Name + ": " + strings.Join(f.Values, "; "))
}
