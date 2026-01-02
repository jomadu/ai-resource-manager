package core

import "testing"

func TestPackageID(t *testing.T) {
	tests := []struct {
		registry string
		name     string
		version  string
		want     string
	}{
		{"test-reg", "test-pkg", "1.0.0", "test-reg/test-pkg@1.0.0"},
		{"my-registry", "my-package", "2.1.3", "my-registry/my-package@2.1.3"},
		{"", "package", "1.0.0", "/package@1.0.0"},
		{"registry", "", "1.0.0", "registry/@1.0.0"},
		{"registry", "package", "", "registry/package@"},
	}

	for _, tt := range tests {
		got := PackageID(tt.registry, tt.name, tt.version)
		if got != tt.want {
			t.Errorf("PackageID(%q, %q, %q) = %q, want %q", 
				tt.registry, tt.name, tt.version, got, tt.want)
		}
	}
}

func TestParsePackageID(t *testing.T) {
	tests := []struct {
		id          string
		wantReg     string
		wantName    string
		wantVersion string
		wantErr     bool
	}{
		{"test-reg/test-pkg@1.0.0", "test-reg", "test-pkg", "1.0.0", false},
		{"my-registry/my-package@2.1.3", "my-registry", "my-package", "2.1.3", false},
		{"registry/package@", "registry", "package", "", false},
		{"invalid-format", "", "", "", true},
		{"registry/package", "", "", "", true},
		{"registry@1.0.0", "", "", "", true},
		{"", "", "", "", true},
	}

	for _, tt := range tests {
		gotReg, gotName, gotVersion, err := ParsePackageID(tt.id)
		
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParsePackageID(%q) expected error, got nil", tt.id)
			}
			continue
		}
		
		if err != nil {
			t.Errorf("ParsePackageID(%q) unexpected error: %v", tt.id, err)
			continue
		}
		
		if gotReg != tt.wantReg || gotName != tt.wantName || gotVersion != tt.wantVersion {
			t.Errorf("ParsePackageID(%q) = (%q, %q, %q), want (%q, %q, %q)", 
				tt.id, gotReg, gotName, gotVersion, tt.wantReg, tt.wantName, tt.wantVersion)
		}
	}
}

func TestPackageIDRoundTrip(t *testing.T) {
	tests := []struct {
		registry string
		name     string
		version  string
	}{
		{"test-reg", "test-pkg", "1.0.0"},
		{"sample-registry", "clean-code-ruleset", "2.1.0"},
		{"my-org", "typescript-rules", "1.5.2"},
	}

	for _, tt := range tests {
		id := PackageID(tt.registry, tt.name, tt.version)
		gotReg, gotName, gotVersion, err := ParsePackageID(id)
		
		if err != nil {
			t.Errorf("Round trip failed with error: %v", err)
			continue
		}
		
		if gotReg != tt.registry || gotName != tt.name || gotVersion != tt.version {
			t.Errorf("Round trip failed: %q/%q@%q -> %q -> %q/%q@%q", 
				tt.registry, tt.name, tt.version, id, gotReg, gotName, gotVersion)
		}
	}
}