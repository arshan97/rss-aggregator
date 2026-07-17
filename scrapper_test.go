package main

import (
	"testing"
	"time"
)

func TestParseFeedDate(t *testing.T) {
	tests := []struct {
		name string
		input string
		want time.Time
	}{
		{
			name:  "RFC1123Z numeric offset",
			input: "Mon, 15 Jun 2026 21:19:29 +0000",
			want:  time.Date(2026, time.June, 15, 21, 19, 29, 0, time.UTC),
		},
		{
			name:  "RFC1123 GMT abbreviation",
			input: "Mon, 15 Jun 2026 21:19:29 GMT",
			want:  time.Date(2026, time.June, 15, 21, 19, 29, 0, time.UTC),
		},
		{
			name:  "RFC3339",
			input: "2026-06-15T21:19:29Z",
			want:  time.Date(2026, time.June, 15, 21, 19, 29, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFeedDate(tt.input)
			if err != nil {
				t.Fatalf("parseFeedDate(%q) returned error: %v", tt.input, err)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("parseFeedDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
