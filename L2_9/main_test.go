package main

import (
	"testing"
)

func TestStringUnpack(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "обычная распаковка",
			input:   "a4bc2d5e",
			want:    "aaaabccddddde",
			wantErr: false,
		},
		{
			name:    "без цифр",
			input:   "abcd",
			want:    "abcd",
			wantErr: false,
		},
		{
			name:    "только цифры",
			input:   "45",
			want:    "",
			wantErr: true,
		},
		{
			name:    "пустая строка",
			input:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "начинается с цифры",
			input:   "1abc",
			want:    "",
			wantErr: true,
		},
		{
			name:    "один символ",
			input:   "a",
			want:    "a",
			wantErr: false,
		},
		{
			name:    "буква и цифра",
			input:   "a3",
			want:    "aaa",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringUnpack(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringUnpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringUnpack() = %q, want %q", got, tt.want)
			}
		})
	}
}
