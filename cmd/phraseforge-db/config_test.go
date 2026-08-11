package main

import (
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		want    config
		wantErr string
	}{
		{
			name: "uses default files",
			want: config{
				dataFile:     defaultDataFile,
				databaseFile: defaultDatabaseFile,
			},
		},
		{
			name: "uses configured files",
			args: []string{"--data-file", "custom.json", "--database-file", "custom.db"},
			want: config{
				dataFile:     "custom.json",
				databaseFile: "custom.db",
			},
		},
		{
			name:    "returns error for unknown flag",
			args:    []string{"--unknown"},
			wantErr: "flag provided but not defined",
		},
		{
			name:    "returns error for positional argument",
			args:    []string{"init"},
			wantErr: "unexpected argument \"init\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfig(tt.args)

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
		})
	}
}
