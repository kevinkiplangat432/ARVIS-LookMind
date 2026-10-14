package commands

import "testing"

func TestValidateIdentityName(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"plain name", "safaricom-billing", "safaricom-billing", false},
		{"trims spaces", "  jane doe  ", "jane doe", false},
		{"empty", "", "", true},
		{"only spaces", "   ", "", true},
		{"only tabs and newlines", "\t\n", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateIdentityName(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}