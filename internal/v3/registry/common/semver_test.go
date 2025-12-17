package common

import (
	"testing"
)

func TestSemverHelper_IsSemverVersion(t *testing.T) {
	helper := NewSemverHelper()

	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{"valid semver", "1.2.3", true},
		{"valid semver with v", "v1.2.3", true},
		{"invalid format", "1.2", false},
		{"branch name", "main", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helper.IsSemverVersion(tt.version)
			if got != tt.want {
				t.Errorf("IsSemverVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemverHelper_SortVersionsBySemver(t *testing.T) {
	helper := NewSemverHelper()

	tests := []struct {
		name     string
		versions []string
		want     []string
	}{
		{
			name:     "mixed versions",
			versions: []string{"1.0.0", "2.1.0", "1.2.0", "1.1.0", "v1.3.0"},
			want:     []string{"2.1.0", "v1.3.0", "1.2.0", "1.1.0", "1.0.0"},
		},
		{
			name:     "with non-semver",
			versions: []string{"1.0.0", "main", "1.1.0", "dev"},
			want:     []string{"1.1.0", "1.0.0", "main", "dev"},
		},
		{
			name:     "empty",
			versions: []string{},
			want:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helper.SortVersionsBySemver(tt.versions)
			if len(got) != len(tt.want) {
				t.Errorf("SortVersionsBySemver() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("SortVersionsBySemver()[%d] = %s, want %s", i, v, tt.want[i])
				}
			}
		})
	}
}
