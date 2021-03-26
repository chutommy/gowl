package gowl_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/chutified/gowl"
)

func TestField_Render(t *testing.T) {
	type fields struct {
		Name   string
		Values []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "single value from",
			fields: fields{
				Name:   "From",
				Values: []string{"John Doe <john.doe@example.com>"},
			},
			want: []byte("From: John Doe <john.doe@example.com>"),
		},
		{
			name: "single value date",
			fields: fields{
				Name:   "Date",
				Values: []string{"Fri, 10 May 2019 16:40:00 +0100"},
			},
			want: []byte("Date: Fri, 10 May 2019 16:40:00 +0100"),
		},
		{
			name: "multiple value content-type",
			fields: fields{
				Name:   "Content-Type",
				Values: []string{"text/plain", "charset=\"UTF-8\""},
			},
			want: []byte("Content-Type: text/plain; charset=\"UTF-8\""),
		},
		{
			name: "multiple value received",
			fields: fields{
				Name:   "Received",
				Values: []string{"by 1010:abc:abcd:0:0:0:0:0 with SMTP id 123456789abcdef", "Sat, 13 Mar 2021 07:00:30 -0800 (PST)"},
			},
			want: []byte("Received: by 1010:abc:abcd:0:0:0:0:0 with SMTP id 123456789abcdef; Sat, 13 Mar 2021 07:00:30 -0800 (PST)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := gowl.Field{
				Name:   tt.fields.Name,
				Values: tt.fields.Values,
			}
			require.Equal(t, string(tt.want), string(f.Render()))
		})
	}
}

func TestHeader_Render(t *testing.T) {
	type fields struct {
		Fields []gowl.Field
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "basic",
			fields: fields{
				Fields: []gowl.Field{
					{Name: "From", Values: []string{"John Doe <john.doe@example.com>"}},
					{Name: "To", Values: []string{"Thomas Smith <thomas.smith@example.com>"}},
					{Name: "Date", Values: []string{"Wed, 8 Mar 2021 12:45:10 +0100"}},
					{Name: "MIME-Version", Values: []string{"1.0"}},
					{Name: "Content-Type", Values: []string{"multipart/alternative", "boundary=\"37a48tbyab7wot468rls798t3y5fcz4t\""}},
				},
			},
			want: []byte(`From: John Doe <john.doe@example.com>
To: Thomas Smith <thomas.smith@example.com>
Date: Wed, 8 Mar 2021 12:45:10 +0100
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="37a48tbyab7wot468rls798t3y5fcz4t"`,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := gowl.Header{
				Fields: tt.fields.Fields,
			}
			require.Equal(t, string(tt.want), string(h.Render()))
		})
	}
}
