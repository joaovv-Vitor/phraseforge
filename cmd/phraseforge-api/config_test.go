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

func TestAPIDataFile(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "uses default data file without configured value",
			value: "",
			want:  defaultDataFile,
		},
		{
			name:  "uses configured data file",
			value: "data/custom-phrases.json",
			want:  "data/custom-phrases.json",
		},
		{
			name:  "uses default data file for whitespace value",
			value: "   ",
			want:  defaultDataFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PHRASEFORGE_DATA_FILE", tt.value)

			if got := apiDataFile(); got != tt.want {
				t.Errorf("apiDataFile() = %q, want %q", got, tt.want)
			}
		})
	}
}
