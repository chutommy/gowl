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
					gowl.NewPart(
						gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
						strings.NewReader("This is a test message."),
						nil,
					),
					gowl.NewPart(
						gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html", "charset=\"UTF-8\""})}),
						strings.NewReader("<div dir=\"ltr\">This is a test message.</div>"),
						nil,
					),
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
					gowl.NewPart(
						gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"multipart/alternative", "boundary=\"0000000000001d296c05be7539bb\""})}),
						nil,
						[]*gowl.Part{
							gowl.NewPart(
								gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
								strings.NewReader("This is a test file."),
								nil,
							),
							gowl.NewPart(
								gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html", "charset=\"UTF-8\""})}),
								strings.NewReader("<div dir=\"ltr\">This is a test file.</div>"),
								nil,
							),
						},
					),
					gowl.NewPart(
						gowl.NewHeader(
							[]*gowl.Field{
								gowl.NewField("Content-Type", []string{"text/plain", "charset=\"US-ASCII\"", "name=\"test.txt\""}),
								gowl.NewField("Content-Disposition", []string{"attachment", "filename=\"test.txt\""}),
								gowl.NewField("Content-Transfer-Encoding", []string{"base64"}),
							},
						),
						strings.NewReader(`VGhpcyBpcyBhIHRlc3QgZmlsZS4K`),
						nil,
					),
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
					gowl.NewPart(

						gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
						strings.NewReader("This is a test message."),
						nil,
					),
					gowl.NewPart(
						gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html", "charset=\"UTF-8\""})}),
						strings.NewReader("<div dir=\"ltr\">This is a test message.</div>"),
						nil,
					),
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
					gowl.NewPart(
						gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
						strings.NewReader("This is a test message."),
						nil,
					),
					gowl.NewPart(
						gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html", "charset=\"UTF-8\""})}),
						strings.NewReader("<div dir=\"ltr\">This is a test message.</div>"),
						nil,
					),
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
					gowl.NewPart(
						gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html", "charset=\"UTF-8\""})}),
						iotest.ErrReader(errors.New("invalid io.Reader")),
						nil,
					),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := gowl.NewPart(
				tt.fields.Header,
				tt.fields.Content,
				tt.fields.Parts,
			)
			got, err := p.Render()
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
