package main

import "testing"

func TestDetectPlatformFor(t *testing.T) {
	tests := []struct {
		name   string
		goos   string
		goarch string
		want   string
	}{
		{
			name:   "darwin arm64",
			goos:   "darwin",
			goarch: "arm64",
			want:   "darwin-arm64",
		},
		{
			name:   "linux amd64",
			goos:   "linux",
			goarch: "amd64",
			want:   "linux-x64-gnu",
		},
		{
			name:   "linux arm64",
			goos:   "linux",
			goarch: "arm64",
			want:   "linux-arm64-gnu",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := detectPlatformFor(tc.goos, tc.goarch)
			if got != tc.want {
				t.Fatalf("detectPlatformFor(%q, %q) = %q, want %q", tc.goos, tc.goarch, got, tc.want)
			}
		})
	}
}
