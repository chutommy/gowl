package gowl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessage_Render(t *testing.T) {
	type fields struct {
		Header *Header
		Root   *Part
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Header: &Header{
					Fields: []*Field{
						{Name: "From", Values: []string{"Johny <john.smith@example.com>"}},
						{Name: "To", Values: []string{"David Doe <david.doe@example.com>"}},
					},
				},
				Root: &Part{
					Header: &Header{
						Fields: []*Field{
							{Name: "Content-Type", Values: []string{"multipart/alternative", `boundary="part_12345"`}},
						},
					},
					Parts: []*Part{
						{
							Header:  &Header{Fields: []*Field{{Name: "Content-Type", Values: []string{"text/plain"}}}},
							Content: strings.NewReader("This is a test message."),
						},
						{
							Header:  &Header{Fields: []*Field{{Name: "Content-Type", Values: []string{"text/html"}}}},
							Content: strings.NewReader(`<div dir="ltr">This is a test message.</div>`),
						},
					},
				},
			},
			want: []byte(`From: Johny <john.smith@example.com>
To: David Doe <david.doe@example.com>Content-Type: multipart/alternative; boundary="part_12345"

--part_12345
Content-Type: text/plain

This is a test message.

--part_12345
Content-Type: text/html

<div dir="ltr">This is a test message.</div>

--part_12345--`),
		},
		{
			name: "header error",
			fields: fields{
				Header: &Header{
					Fields: []*Field{
						{Name: "From", Values: []string{"Johny <john.smith@example.com>"}},
						{Name: "To", Values: nil},
					},
				},
				Root: &Part{
					Header: &Header{
						Fields: []*Field{
							{Name: "Content-Type", Values: []string{"multipart/alternative", `boundary="part_12345"`}},
						},
					},
					Parts: []*Part{
						{
							Header:  &Header{Fields: []*Field{{Name: "Content-Type", Values: []string{"text/plain"}}}},
							Content: strings.NewReader("This is a test message."),
						},
						{
							Header:  &Header{Fields: []*Field{{Name: "Content-Type", Values: []string{"text/html"}}}},
							Content: strings.NewReader(`<div dir="ltr">This is a test message.</div>`),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "root error",
			fields: fields{
				Header: &Header{
					Fields: []*Field{
						{Name: "From", Values: []string{"Johny <john.smith@example.com>"}},
						{Name: "To", Values: []string{"David Doe <david.doe@example.com>"}},
					},
				},
				Root: &Part{
					Header: &Header{
						Fields: []*Field{
							{Name: "Content-Type", Values: []string{"multipart/alternative"}},
						},
					},
					Parts: []*Part{
						{
							Header:  &Header{Fields: []*Field{{Name: "Content-Type", Values: []string{"text/plain"}}}},
							Content: strings.NewReader("This is a test message."),
						},
						{
							Header:  &Header{Fields: []*Field{{Name: "Content-Type", Values: []string{"text/html"}}}},
							Content: strings.NewReader(`<div dir="ltr">This is a test message.</div>`),
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Header: tt.fields.Header,
				Root:   tt.fields.Root,
			}
			got, err := m.Render()
			if tt.wantErr {
				require.Nil(t, got)
				require.Error(t, err)
			} else {
				require.Equal(t, string(tt.want), string(got))
				require.NoError(t, err)
			}
		})
	}
}
