package gowl_test

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"

	"github.com/chutified/gowl"
)

func TestPart_Render(t *testing.T) {
	type fields struct {
		Header  *gowl.Header
		Content io.Reader
		Parts   []*gowl.Part
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "plain and html text",
			fields: fields{
				Header: gowl.NewHeader(
					[]*gowl.Field{
						gowl.NewField("Content-Type", []string{"multipart/alternative", "boundary=\"0000000000009c8ab105be4e2cc3\""}),
					},
				),
				Parts: []*gowl.Part{
					{
						Header:  gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
						Content: strings.NewReader("This is a test message."),
					},
					{
						Header:  gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html", "charset=\"UTF-8\""})}),
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
				Header: gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"multipart/mixed", "boundary=\"0000000000001d296f05be7539bd\""})}),
				Parts: []*gowl.Part{
					{
						Header: gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"multipart/alternative", "boundary=\"0000000000001d296c05be7539bb\""})}),
						Parts: []*gowl.Part{
							{
								Header:  gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
								Content: strings.NewReader("This is a test file."),
							},
							{
								Header:  gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html", "charset=\"UTF-8\""})}),
								Content: strings.NewReader("<div dir=\"ltr\">This is a test file.</div>"),
							},
						},
					},
					{
						Header: gowl.NewHeader(
							[]*gowl.Field{
								gowl.NewField("Content-Type", []string{"text/plain", "charset=\"US-ASCII\"", "name=\"test.txt\""}),
								gowl.NewField("Content-Disposition", []string{"attachment", "filename=\"test.txt\""}),
								gowl.NewField("Content-Transfer-Encoding", []string{"base64"}),
							},
						),
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
		{
			name: "invalid header",
			fields: fields{
				Header: gowl.NewHeader(
					[]*gowl.Field{
						gowl.NewField("Content-Type", nil),
					},
				),
				Parts: []*gowl.Part{
					{
						Header:  gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
						Content: strings.NewReader("This is a test message."),
					},
					{
						Header:  gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html", "charset=\"UTF-8\""})}),
						Content: strings.NewReader("<div dir=\"ltr\">This is a test message.</div>"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "plain text",
			fields: fields{
				Header: gowl.NewHeader(
					[]*gowl.Field{
						gowl.NewField("Content-Type", []string{"text/plain"}),
					},
				),
				Content: strings.NewReader("This is a test message."),
			},
			want: []byte(`Content-Type: text/plain

This is a test message.`,
			),
		},
		{
			name: "invalid content",
			fields: fields{
				Header: gowl.NewHeader(
					[]*gowl.Field{
						gowl.NewField("Content-Type", []string{"text/plain"}),
					},
				),
				Content: iotest.ErrReader(errors.New("invalid io.Reader")),
			},
			wantErr: true,
		},
		{
			name: "missing boundary in header",
			fields: fields{
				Header: gowl.NewHeader(
					[]*gowl.Field{
						gowl.NewField("Content-Type", []string{"multipart/alternative"}),
					},
				),
				Parts: []*gowl.Part{
					{
						Header:  gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
						Content: strings.NewReader("This is a test message."),
					},
					{
						Header:  gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html", "charset=\"UTF-8\""})}),
						Content: strings.NewReader("<div dir=\"ltr\">This is a test message.</div>"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid reader in parts",
			fields: fields{
				Header: gowl.NewHeader(
					[]*gowl.Field{
						gowl.NewField("Content-Type", []string{"multipart/alternative", "boundary=\"0000000000009c8ab105be4e2cc3\""}),
					},
				),
				Parts: []*gowl.Part{
					{
						Header:  gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html", "charset=\"UTF-8\""})}),
						Content: iotest.ErrReader(errors.New("invalid io.Reader")),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// construct Part
			p := &gowl.Part{
				Header:  tt.fields.Header,
				Content: tt.fields.Content,
				Parts:   tt.fields.Parts,
			}

			// run Render
			got, err := p.Render()

			// check returned values
			if tt.wantErr {
				require.Nil(t, got)
				require.Error(t, err)
			} else {
				require.Equal(t, string(tt.want), string(got))
				require.Nil(t, err)
			}
		})
	}
}
