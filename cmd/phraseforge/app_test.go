package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr string
	}{
		{
			name: "generates a programming phrase",
			args: []string{"generate"},
		},
		{
			name: "lists categories in configured order",
			args: []string{"categories"},
			want: "programming\nstudy\n",
		},
		{
			name:    "returns error without command",
			wantErr: "command is required",
		},
		{
			name:    "returns error for unknown command",
			args:    []string{"invalid"},
			wantErr: "unknown command \"invalid\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			err := run(tt.args, &output)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("run() error = nil, want an error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("run() error = %q, want it to contain %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("run() unexpected error: %v", err)
			}

			if tt.name == "generates a programming phrase" {
				if got := output.String(); got == "" || !strings.HasSuffix(got, ".\n") {
					t.Errorf("run() output = %q, want a non-empty phrase ending in a period", got)
				}
				return
			}

			if got := output.String(); got != tt.want {
				t.Errorf("run() output = %q, want %q", got, tt.want)
			}
		})
	}
}
