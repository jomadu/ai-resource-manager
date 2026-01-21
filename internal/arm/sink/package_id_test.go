package sink

import "testing"

func TestPkgKey(t *testing.T) {
	tests := []struct {
		registry string
		name     string
		version  string
		want     string
	}{
		{"reg", "pkg", "1.0.0", "reg/pkg@1.0.0"},
		{"my-reg", "my-pkg", "2.1.3", "my-reg/my-pkg@2.1.3"},
	}

	for _, tt := range tests {
		got := pkgKey(tt.registry, tt.name, tt.version)
		if got != tt.want {
			t.Errorf("pkgKey(%q, %q, %q) = %q, want %q", 
				tt.registry, tt.name, tt.version, got, tt.want)
		}
	}
}

func TestParsePkgKey(t *testing.T) {
	tests := []struct {
		key         string
		wantReg     string
		wantName    string
		wantVersion string
		wantErr     bool
	}{
		{"reg/pkg@1.0.0", "reg", "pkg", "1.0.0", false},
		{"my-reg/my-pkg@2.1.3", "my-reg", "my-pkg", "2.1.3", false},
		{"invalid", "", "", "", true},
		{"reg/pkg", "", "", "", true},
		{"reg@1.0.0", "", "", "", true},
	}

	for _, tt := range tests {
		gotReg, gotName, gotVersion, err := parsePkgKey(tt.key)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parsePkgKey(%q) expected error, got nil", tt.key)
			}
		} else {
			if err != nil {
				t.Errorf("parsePkgKey(%q) unexpected error: %v", tt.key, err)
			}
			if gotReg != tt.wantReg || gotName != tt.wantName || gotVersion != tt.wantVersion {
				t.Errorf("parsePkgKey(%q) = (%q, %q, %q), want (%q, %q, %q)", 
					tt.key, gotReg, gotName, gotVersion, tt.wantReg, tt.wantName, tt.wantVersion)
			}
		}
	}
}

func TestPkgKeyRoundTrip(t *testing.T) {
	tests := []struct {
		registry string
		name     string
		version  string
	}{
		{"reg", "pkg", "1.0.0"},
		{"my-reg", "my-pkg", "2.1.3"},
	}

	for _, tt := range tests {
		key := pkgKey(tt.registry, tt.name, tt.version)
		gotReg, gotName, gotVersion, err := parsePkgKey(key)
		if err != nil {
			t.Errorf("Round trip failed: %v", err)
		}
		if gotReg != tt.registry || gotName != tt.name || gotVersion != tt.version {
			t.Errorf("Round trip mismatch: got (%q, %q, %q), want (%q, %q, %q)",
				gotReg, gotName, gotVersion, tt.registry, tt.name, tt.version)
		}
	}
}
