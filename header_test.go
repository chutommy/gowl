package gowl_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/chutified/gowl"
)

func TestField_Param(t *testing.T) {
	type fields struct {
		Name   string
		Values []string
	}
	type args struct {
		param string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "content type boundary",
			fields: fields{
				Name:   "Content-Type",
				Values: []string{"multipart/mixed", `boundary="part_1234567890"`},
			},
			args: args{
				param: "boundary",
			},
			want: []byte("part_1234567890"),
		},
		{
			name: "content type value",
			fields: fields{
				Name:   "Content-Type",
				Values: []string{"text/plain", `charset="UTF-8"`},
			},
			args: args{
				param: "charset",
			},
			want: []byte("UTF-8"),
		},
		{
			name: "content type no value",
			fields: fields{
				Name:   "Content-Type",
				Values: []string{"text/plain", `charset="UTF-8"`},
			},
			args: args{
				param: "boundary",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &gowl.Field{
				Name:   tt.fields.Name,
				Values: tt.fields.Values,
			}
			got := f.Param(tt.args.param)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestField_Render(t *testing.T) {
	type fields struct {
		Name   string
		Values []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr error
	}{
		{
			name: "single value from",
			fields: fields{
				Name:   "From",
				Values: []string{"John Doe <john.doe@example.com>"},
			},
			want:    []byte("From: John Doe <john.doe@example.com>"),
			wantErr: nil,
		},
		{
			name: "single value date",
			fields: fields{
				Name:   "Date",
				Values: []string{"Fri, 10 May 2019 16:40:00 +0100"},
			},
			want:    []byte("Date: Fri, 10 May 2019 16:40:00 +0100"),
			wantErr: nil,
		},
		{
			name: "multiple value content-type",
			fields: fields{
				Name:   "Content-Type",
				Values: []string{"text/plain", "charset=\"UTF-8\""},
			},
			want:    []byte("Content-Type: text/plain; charset=\"UTF-8\""),
			wantErr: nil,
		},
		{
			name: "multiple value received",
			fields: fields{
				Name:   "Received",
				Values: []string{"by 1010:abc:abcd:0:0:0:0:0 with SMTP id 123456789abcdef", "Sat, 13 Mar 2021 07:00:30 -0800 (PST)"},
			},
			want:    []byte("Received: by 1010:abc:abcd:0:0:0:0:0 with SMTP id 123456789abcdef; Sat, 13 Mar 2021 07:00:30 -0800 (PST)"),
			wantErr: nil,
		},
		{
			name:    "zero values",
			fields:  fields{},
			wantErr: gowl.ErrNoValues,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// construct Field
			f := &gowl.Field{
				Name:   tt.fields.Name,
				Values: tt.fields.Values,
			}

			// run Render
			got, err := f.Render()

			// check returned values
			if tt.wantErr != nil {
				require.Nil(t, got)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.Equal(t, string(tt.want), string(got))
				require.Nil(t, err)
			}
		})
	}
}

func TestHeader_Render(t *testing.T) {
	type fields struct {
		Fields []*gowl.Field
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				Fields: []*gowl.Field{
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
			wantErr: nil,
		},
		{
			name: "no boundary",
			fields: fields{
				Fields: []*gowl.Field{
					{Name: "From", Values: []string{"John Doe <john.doe@example.com>"}},
					{Name: "To", Values: []string{"Thomas Smith <thomas.smith@example.com>"}},
					{Name: "MIME-Version", Values: []string{"1.0"}},
					{Name: "Content-Type", Values: nil},
				},
			},
			wantErr: gowl.ErrNoValues,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// construct Header
			h := &gowl.Header{
				Fields: tt.fields.Fields,
			}

			// run Render
			got, err := h.Render()

			// check returned values
			if tt.wantErr != nil {
				require.Nil(t, got)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.Equal(t, string(tt.want), string(got))
				require.Nil(t, err)
			}
		})
	}
}

func TestHeader_Boundary(t *testing.T) {
	type fields struct {
		Fields []*gowl.Field
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				Fields: []*gowl.Field{
					{Name: "From", Values: []string{"<david.smith@example.com>"}},
					{Name: "To", Values: []string{"<john.doe@example.com>"}},
					{Name: "Content-Type", Values: []string{"multipart/alternative", `boundary="part_12345"`}},
				},
			},
			want:    []byte("part_12345"),
			wantErr: nil,
		},
		{
			name: "no boundary",
			fields: fields{
				Fields: []*gowl.Field{
					{Name: "From", Values: []string{"<david.smith@example.com>"}},
					{Name: "To", Values: []string{"<john.doe@example.com>"}},
					{Name: "Content-Type", Values: []string{"text/plain", `charset="UTF-8"`}},
				},
			},
			wantErr: gowl.ErrNoBoundary,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &gowl.Header{
				Fields: tt.fields.Fields,
			}
			got, err := h.Boundary()
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
			} else {
				require.Nil(t, err)
				require.Equal(t, string(tt.want), string(got))
			}
		})
	}
}
