package gowl_test

import (
	"strings"
	"testing"

	"github.com/chutommy/gowl"
	"github.com/stretchr/testify/require"
)

func TestMessage_Reset(t *testing.T) {
	t.Parallel()

	h := gowl.NewHeader([]*gowl.Field{
		gowl.NewField("From", []string{"John Smith <john.smith@example.com>"}),
		gowl.NewField("To", []string{"<thomas.harold@example.com>"}),
	})
	rp := gowl.NewPart(
		gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
		strings.NewReader("This is a test message."),
		nil,
	)

	msg := gowl.NewMessage(h, rp)
	msg.Reset()
	require.Equal(t, &gowl.Message{}, msg)
}

func TestMessage_Header(t *testing.T) {
	t.Parallel()

	h := gowl.NewHeader([]*gowl.Field{
		gowl.NewField("From", []string{"John Smith <john.smith@example.com>"}),
		gowl.NewField("To", []string{"<thomas.harold@example.com>"}),
	})

	msg := gowl.NewMessage(h, nil)
	got := msg.Header()
	require.Equal(t, h, got)
}

func TestMessage_RootPart(t *testing.T) {
	t.Parallel()

	rp := gowl.NewPart(
		gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
		strings.NewReader("This is a test message."),
		nil,
	)

	msg := gowl.NewMessage(nil, rp)
	got := msg.RootPart()
	require.Equal(t, rp, got)
}

func TestMessage_SetHeader(t *testing.T) {
	t.Parallel()

	h := gowl.NewHeader([]*gowl.Field{
		gowl.NewField("From", []string{"John Smith <john.smith@example.com>"}),
		gowl.NewField("To", []string{"<thomas.harold@example.com>"}),
	})
	h2 := gowl.NewHeader([]*gowl.Field{
		gowl.NewField("From", []string{"David Doe<david.doe@example.com>"}),
		gowl.NewField("To", []string{"<marcus.white@example.com>"}),
	})

	msg := gowl.NewMessage(h, nil)
	msg.SetHeader(h2)
	got := msg.Header()
	require.Equal(t, h2, got)
}

func TestMessage_SetRootPart(t *testing.T) {
	t.Parallel()

	rp := gowl.NewPart(
		gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain", "charset=\"UTF-8\""})}),
		strings.NewReader("This is a test message."),
		nil,
	)
	rp2 := gowl.NewPart(
		gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html"})}),
		strings.NewReader(`<div dir="ltr">This is a test message.<dir>`),
		nil,
	)

	msg := gowl.NewMessage(nil, rp)
	msg.SetRootPart(rp2)
	got := msg.RootPart()
	require.Equal(t, rp2, got)
}

func TestMessage_Render(t *testing.T) {
	t.Parallel()

	type fields struct {
		Header *gowl.Header
		Root   *gowl.Part
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
				Header: gowl.NewHeader(
					[]*gowl.Field{
						gowl.NewField("From", []string{"Johny <john.smith@example.com>"}),
						gowl.NewField("To", []string{"David Doe <david.doe@example.com>"}),
					},
				),
				Root: gowl.NewPart(
					gowl.NewHeader(
						[]*gowl.Field{
							gowl.NewField("Content-Type", []string{"multipart/alternative", `boundary="part_12345"`}),
						},
					),
					nil,
					[]*gowl.Part{
						gowl.NewPart(
							gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain"})}),
							strings.NewReader("This is a test message."),
							nil,
						),
						gowl.NewPart(
							gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html"})}),
							strings.NewReader(`<div dir="ltr">This is a test message.</div>`),
							nil,
						),
					},
				),
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
				Header: gowl.NewHeader(
					[]*gowl.Field{
						gowl.NewField("From", []string{"Johny <john.smith@example.com>"}),
						gowl.NewField("To", nil),
					},
				),
				Root: gowl.NewPart(
					gowl.NewHeader(
						[]*gowl.Field{
							gowl.NewField("Content-Type", []string{"multipart/alternative", `boundary="part_12345"`}),
						},
					),
					nil,
					[]*gowl.Part{
						gowl.NewPart(
							gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain"})}),
							strings.NewReader("This is a test message."),
							nil,
						),
						gowl.NewPart(
							gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html"})}),
							strings.NewReader(`<div dir="ltr">This is a test message.</div>`),
							nil,
						),
					},
				),
			},
			wantErr: true,
		},
		{
			name: "root error",
			fields: fields{
				Header: gowl.NewHeader(
					[]*gowl.Field{
						gowl.NewField("From", []string{"Johny <john.smith@example.com>"}),
						gowl.NewField("To", []string{"David Doe <david.doe@example.com>"}),
					},
				),
				Root: gowl.NewPart(
					gowl.NewHeader(
						[]*gowl.Field{
							gowl.NewField("Content-Type", []string{"multipart/alternative"}),
						},
					),
					nil,
					[]*gowl.Part{
						gowl.NewPart(
							gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/plain"})}),
							strings.NewReader("This is a test message."),
							nil,
						),
						gowl.NewPart(
							gowl.NewHeader([]*gowl.Field{gowl.NewField("Content-Type", []string{"text/html"})}),
							strings.NewReader(`<div dir="ltr">This is a test message.</div>`),
							nil,
						),
					},
				),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := gowl.NewMessage(tt.fields.Header, tt.fields.Root)
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
