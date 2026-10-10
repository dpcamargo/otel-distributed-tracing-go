package main

import "testing"

func TestValidateZip(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"valid 8 digits", "01001000", "01001000", false},
		{"valid with dash", "01001-000", "01001000", false},
		{"valid with spaces", "01001 000", "01001000", false},
		{"too short", "0100100", "", true},
		{"too long", "010010000", "", true},
		{"empty", "", "", true},
		{"non-numeric", "abcdefgh", "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := validateZip(c.input)
			if c.wantErr {
				if err == nil {
					t.Fatalf("validateZip(%q): expected error, got nil (result %q)", c.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("validateZip(%q): unexpected error: %v", c.input, err)
			}
			if got != c.want {
				t.Fatalf("validateZip(%q) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}
