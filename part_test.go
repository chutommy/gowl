package gowl

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPart_Render(t *testing.T) {
	type fields struct {
		Header  Header
		Content io.Reader
		Parts   []*Part
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "plain and html text",
			fields: fields{
				Header: Header{
					Fields: []Field{
						{Name: "Content-Type", Values: []string{"multipart/alternative", "boundary=\"0000000000009c8ab105be4e2cc3\""}},
					},
				},
				Parts: []*Part{
					{
						Header:  Header{Fields: []Field{{Name: "Content-Type", Values: []string{"text/plain", "charset=\"UTF-8\""}}}},
						Content: strings.NewReader("This is a test message."),
					},
					{
						Header:  Header{Fields: []Field{{Name: "Content-Type", Values: []string{"text/html", "charset=\"UTF-8\""}}}},
						Content: strings.NewReader("<div dir=\"ltr\">This is a test message.</div>"),
					},
				},
			},
			want: []byte(`Content-Type: multipart/alternative; boundary="0000000000009c8ab105be4e2cc3"

--0000000000009c8ab105be4e2cc3
Content-Type: text/plain; charset="UTF-8"

This is a test message.

--0000000000009c8ab105be4e2cc3
Content-Type: text/html; charset="UTF-8"

<div dir="ltr">This is a test message.</div>

--0000000000009c8ab105be4e2cc3--`,
			),
		},
		{
			name: "text with an attachment",
			fields: fields{
				Header: Header{Fields: []Field{{Name: "Content-Type", Values: []string{"multipart/mixed", "boundary=\"0000000000001d296f05be7539bd\""}}}},
				Parts: []*Part{
					{
						Header: Header{Fields: []Field{{Name: "Content-Type", Values: []string{"multipart/alternative", "boundary=\"0000000000001d296c05be7539bb\""}}}},
						Parts: []*Part{
							{
								Header:  Header{Fields: []Field{{Name: "Content-Type", Values: []string{"text/plain", "charset=\"UTF-8\""}}}},
								Content: strings.NewReader("This is a test file."),
							},
							{
								Header:  Header{Fields: []Field{{Name: "Content-Type", Values: []string{"text/html", "charset=\"UTF-8\""}}}},
								Content: strings.NewReader("<div dir=\"ltr\">This is a test file.</div>"),
							},
						},
					},
					{
						Header: Header{
							Fields: []Field{
								{Name: "Content-Type", Values: []string{"text/plain", "charset=\"US-ASCII\"", "name=\"test.txt\""}},
								{Name: "Content-Disposition", Values: []string{"attachment", "filename=\"test.txt\""}},
								{Name: "Content-Transfer-Encoding", Values: []string{"base64"}},
							},
						},
						Content: strings.NewReader(`VGhpcyBpcyBhIHRlc3QgZmlsZS4K`),
					},
				},
			},
			want: []byte(`Content-Type: multipart/mixed; boundary="0000000000001d296f05be7539bd"

--0000000000001d296f05be7539bd
Content-Type: multipart/alternative; boundary="0000000000001d296c05be7539bb"

--0000000000001d296c05be7539bb
Content-Type: text/plain; charset="UTF-8"

This is a test file.

--0000000000001d296c05be7539bb
Content-Type: text/html; charset="UTF-8"

<div dir="ltr">This is a test file.</div>

--0000000000001d296c05be7539bb--

--0000000000001d296f05be7539bd
Content-Type: text/plain; charset="US-ASCII"; name="test.txt"
Content-Disposition: attachment; filename="test.txt"
Content-Transfer-Encoding: base64

VGhpcyBpcyBhIHRlc3QgZmlsZS4K

--0000000000001d296f05be7539bd--`,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Part{
				Header:  tt.fields.Header,
				Content: tt.fields.Content,
				Parts:   tt.fields.Parts,
			}
			require.Equal(t, string(tt.want), string(p.Render()))
		})
	}
}
