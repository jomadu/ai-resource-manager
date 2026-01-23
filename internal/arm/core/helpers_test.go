package core

import "testing"

func TestPackageKey(t *testing.T) {
	tests := []struct {
		registry string
		pkg      string
		want     string
	}{
		{"registry1", "package1", "registry1/package1"},
		{"my-registry", "my-package", "my-registry/my-package"},
		{"", "package", "/package"},
		{"registry", "", "registry/"},
	}

	for _, tt := range tests {
		got := PackageKey(tt.registry, tt.pkg)
		if got != tt.want {
			t.Errorf("PackageKey(%q, %q) = %q, want %q", tt.registry, tt.pkg, got, tt.want)
		}
	}
}

func TestParsePackageKey(t *testing.T) {
	tests := []struct {
		key     string
		wantReg string
		wantPkg string
	}{
		{"registry1/package1", "registry1", "package1"},
		{"my-registry/my-package", "my-registry", "my-package"},
		{"registry/package/with/slashes", "registry", "package/with/slashes"},
		{"invalid", "", ""},
		{"", "", ""},
		{"/package", "", "package"},
		{"registry/", "registry", ""},
	}

	for _, tt := range tests {
		gotReg, gotPkg := ParsePackageKey(tt.key)
		if gotReg != tt.wantReg || gotPkg != tt.wantPkg {
			t.Errorf("ParsePackageKey(%q) = (%q, %q), want (%q, %q)",
				tt.key, gotReg, gotPkg, tt.wantReg, tt.wantPkg)
		}
	}
}

func TestPackageKeyRoundTrip(t *testing.T) {
	tests := []struct {
		registry string
		pkg      string
	}{
		{"registry1", "package1"},
		{"sample-registry", "clean-code-ruleset"},
		{"my-org", "typescript-rules"},
	}

	for _, tt := range tests {
		key := PackageKey(tt.registry, tt.pkg)
		gotReg, gotPkg := ParsePackageKey(key)

		if gotReg != tt.registry || gotPkg != tt.pkg {
			t.Errorf("Round trip failed: %q/%q -> %q -> %q/%q",
				tt.registry, tt.pkg, key, gotReg, gotPkg)
		}
	}
}
