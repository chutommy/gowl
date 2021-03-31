package gowl

import (
	"bytes"
	"errors"
	"strings"
)

// Error codes returned by failures to render to SMTP.
var (
	ErrNoValues   = errors.New("the attribute values of Field is empty")
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
	for _, f := range h.Fields {
		if f.name == "Content-Type" {
			if v := f.Param("boundary"); v != nil {
				return v, nil
			}
		}
	}

	return nil, ErrNoBoundary
}

// Field represents a single SMTP Header field.
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

// Name returns the Field name.
func (f *Field) Name() string {
	return f.name
}

// Name returns the Field values.
func (f *Field) Values() []string {
	return f.values
}

// SetName rewrites the Field name value.
func (f *Field) SetName(name string) {
	f.name = name
}

// SetName rewrites the Field values value.
func (f *Field) SetValues(values []string) {
	f.values = values
}

// AddValue appends given value to the end of the Field values.
func (f *Field) AddValue(value string) {
	f.values = append(f.values, value)
}

// Param returns the value of the parameter param of the Field f.
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
// SMTP Header Field. The Field's values are separated with semicolons.
func (f *Field) Render() ([]byte, error) {
	if len(f.values) == 0 {
		return nil, ErrNoValues
	}

	return []byte(f.name + ": " + strings.Join(f.values, "; ")), nil
}
