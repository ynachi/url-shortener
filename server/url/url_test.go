package server

import (
	"reflect"
	"testing"
)

func TestToBase58(t *testing.T) {
	// base64 numeric
	in1 := "4404500914867889436"
	in2 := "14629859553178350951"
	tests := []struct {
		name string
		give string
		want string
	}{
		{
			name: "working_1",
			give: in1,
			want: "BDzDPKK56kj",
		},
		{
			name: "working_2",
			give: in2,
			want: "axeUhbEvZgz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := toBase58([]byte(tt.give))
			if got != tt.want {
				const msg = `server.toBase58 : wanted %v but got %v.`
				t.Fatalf(msg, tt.want, got)
			}
		})
	}
}

func TestToSHAR256(t *testing.T) {
	in1 := "http://www.golang.com"
	in2 := "http://www.yoa.com"
	tests := []struct {
		name string
		give string
		want []byte
	}{
		{
			name: "working_1",
			give: in1,
			want: []byte{198, 170, 8, 60, 67, 95, 120, 29, 234, 238, 13, 47, 252, 185, 246, 150, 184, 165, 239, 142, 167, 177, 207, 130, 61, 31, 238, 62, 251, 149, 101, 28},
		},
		{
			name: "working_2",
			give: in2,
			want: []byte{194, 156, 248, 133, 242, 213, 231, 157, 188, 94, 208, 76, 111, 224, 123, 94, 14, 152, 214, 159, 108, 90, 16, 140, 203, 7, 179, 192, 85, 12, 245, 103},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSHA256(tt.give)
			if !reflect.DeepEqual(got, tt.want) {
				const msg = `server.toBase58 : wanted %v but got %v.`
				t.Fatalf(msg, tt.want, got)
			}
		})
	}
}
