package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     []string
		want     config
		wantArgs []string
		wantErr  string
	}{
		{
			name:     "uses default data file",
			args:     []string{"categories"},
			want:     config{dataFile: defaultDataFile},
			wantArgs: []string{"categories"},
		},
		{
			name:     "uses configured data file",
			args:     []string{"--data-file", "custom.json", "generate", "--count", "3"},
			want:     config{dataFile: "custom.json"},
			wantArgs: []string{"generate", "--count", "3"},
		},
		{
			name:    "returns error for unknown application flag",
			args:    []string{"--unknown", "categories"},
			wantErr: "flag provided but not defined",
		},
		{
			name:    "returns error when data file value is missing",
			args:    []string{"--data-file"},
			wantErr: "flag needs an argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotArgs, err := parseConfig(tt.args)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("parseConfig() error = nil, want an error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("parseConfig() error = %q, want it to contain %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("parseConfig() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("parseConfig() config = %#v, want %#v", got, tt.want)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("parseConfig() remaining arguments = %#v, want %#v", gotArgs, tt.wantArgs)
			}
		})
	}
}
