package utils

import (
	"testing"
)

func TestUUID(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UUID()
			if (err != nil) != tt.wantErr {
				t.Errorf("UUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got2, _ := UUID()
			if got == got2 {
				t.Errorf("UUID() should be unique, got = %s, secondUUID = %s", got, got2)
			}
		})
	}
}

func TestContainsUUID(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "uuid",
			want: true,
			args: args{
				input: "11BB4040-7162-AC2D-8DD8-C98FFC7D871D",
			},
		},
		{
			name: "not uuid",
			want: false,
			args: args{
				input: "not a uid",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsUUID(tt.args.input); got != tt.want {
				t.Errorf("ContainsUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractUUID(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "success",
			url:  "http://127.0.0.1:8080/api/health/checks/3b447fdf-d2e9-42bd-adcf-77d147b8b4dc",
			want: "3b447fdf-d2e9-42bd-adcf-77d147b8b4dc",
		},
		{
			name: "invalid uuid",
			url:  "http://127.0.0.1:8080/api/health/checks/3b447fdf-d2e9-42bd-adcf-77",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractUUID(tt.url); got != tt.want {
				t.Errorf("ExtractUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
