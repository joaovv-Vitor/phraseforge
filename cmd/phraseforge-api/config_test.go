package main

import "testing"

func TestAPIAddress(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "uses default address without configured value",
			value: "",
			want:  defaultAPIAddress,
		},
		{
			name:  "uses configured address",
			value: ":9090",
			want:  ":9090",
		},
		{
			name:  "uses default address for whitespace value",
			value: "   ",
			want:  defaultAPIAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PHRASEFORGE_API_ADDR", tt.value)

			if got := apiAddress(); got != tt.want {
				t.Errorf("apiAddress() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIDatabaseFile(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "uses default database file without configured value",
			value: "",
			want:  defaultDatabaseFile,
		},
		{
			name:  "uses configured database file",
			value: "data/custom-phraseforge.db",
			want:  "data/custom-phraseforge.db",
		},
		{
			name:  "uses default database file for whitespace value",
			value: "   ",
			want:  defaultDatabaseFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PHRASEFORGE_DATABASE_FILE", tt.value)

			if got := apiDatabaseFile(); got != tt.want {
				t.Errorf("apiDatabaseFile() = %q, want %q", got, tt.want)
			}
		})
	}
}
