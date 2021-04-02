package gowl

import (
	"bytes"
	"errors"
	"strings"
)

// Error codes returned by failures to render an SMTP data.
var (
	ErrNoValues   = errors.New("the attribute values of Field is empty")
	ErrNoBoundary = errors.New("the Header has no Content-Type field with boundary parameter")
)

// Header represents an SMTP Header.
type Header struct {
	fields []*Field
}

// NewHeader is a constructor for the Header.
func NewHeader(fields []*Field) *Header {
	return &Header{
		fields: fields,
	}
}

// Reset resets the value of the Header but it keeps its instance (pointer).
func (h *Header) Reset() {
	*h = Header{}
}

// Fields returns a list of fields in the Header.
func (h *Header) Fields() []*Field {
	return h.fields
}

// AddField appends a given field to the end of the fields of the Header.
func (h *Header) AddField(field *Field) {
	h.fields = append(h.fields, field)
}

// RemoveField removes a field with a given name in the fields of the Header.
func (h *Header) RemoveField(name string) {
	for i, f := range h.fields {
		if f.name == name {
			h.fields = append(h.fields[:i], h.fields[i+1:]...)

			break
		}
	}
}

// Render renders the Header fields and returns them in bytes.
// It renders each field on its own line.
func (h *Header) Render() ([]byte, error) {
	fs := make([][]byte, len(h.fields))

	var err error
	for i, f := range h.fields {
		if fs[i], err = f.Render(); err != nil {
			return nil, err
		}
	}

	return bytes.Join(fs, []byte{'\n'}), nil
}

// Boundary queries the fields of the Header and tries to find a boundary of the Content-Type.
// If there's no boundary parameter inside its Content-Type ErrNoBoundary is returned.
func (h *Header) Boundary() ([]byte, error) {
	for _, f := range h.fields {
		if f.name == "Content-Type" {
			if v := f.Param("boundary"); v != nil {
				return v, nil
			}
		}
	}

	return nil, ErrNoBoundary
}

// Field represents a single SMTP field of the Header.
type Field struct {
	name   string
	values []string
}

// NewField is a constructor of a Field.
func NewField(name string, values []string) *Field {
	return &Field{
		name:   name,
		values: values,
	}
}

// Reset resets the value of the Field but it keeps its instance (pointer).
func (f *Field) Reset() {
	*f = Field{}
}

// Name returns the name of the Field.
func (f *Field) Name() string {
	return f.name
}

// Values returns the values of the Field.
func (f *Field) Values() []string {
	return f.values
}

// SetName rewrites the name of the value of the Field.
func (f *Field) SetName(name string) {
	f.name = name
}

// SetValues rewrites the values value of the Field.
func (f *Field) SetValues(values []string) {
	f.values = values
}

// AddValue appends given value to the end of the values of the Field.
func (f *Field) AddValue(value string) {
	f.values = append(f.values, value)
}

// Param returns the value of the parameter param of the Field.
func (f *Field) Param(param string) []byte {
	param += "="
	for _, v := range f.values {
		if strings.Contains(v, param) {
			return []byte(strings.Trim(strings.TrimPrefix(v, param), "\""))
		}
	}

	return nil
}

// Render renders the content of the field into bytes. It returns formatted
// SMTP Field of the Header. The values of the Field are separated by semicolons.
func (f *Field) Render() ([]byte, error) {
	if len(f.values) == 0 {
		return nil, ErrNoValues
	}

	return []byte(f.name + ": " + strings.Join(f.values, "; ")), nil
}
