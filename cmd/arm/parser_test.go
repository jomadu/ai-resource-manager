package main

import (
	"reflect"
	"testing"
)

func TestParsePackageArg(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    PackageRef
		wantErr bool
	}{
		{
			name: "registry and package only",
			arg:  "registry/package",
			want: PackageRef{Registry: "registry", Name: "package", Version: ""},
		},
		{
			name: "with version",
			arg:  "registry/package@1.0.0",
			want: PackageRef{Registry: "registry", Name: "package", Version: "1.0.0"},
		},
		{
			name: "with branch version",
			arg:  "registry/package@main",
			want: PackageRef{Registry: "registry", Name: "package", Version: "main"},
		},
		{
			name:    "empty arg",
			arg:     "",
			wantErr: true,
		},
		{
			name:    "missing registry",
			arg:     "package",
			wantErr: true,
		},
		{
			name:    "empty registry",
			arg:     "/package",
			wantErr: true,
		},
		{
			name:    "empty package",
			arg:     "registry/",
			wantErr: true,
		},
		{
			name:    "empty package with version",
			arg:     "registry/@1.0.0",
			wantErr: true,
		},
		{
			name: "complex registry name",
			arg:  "github.com/user/package@v1.0.0",
			want: PackageRef{Registry: "github.com", Name: "user/package", Version: "v1.0.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePackageArg(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePackageArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePackageArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePackageArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []PackageRef
		wantErr bool
	}{
		{
			name: "single arg",
			args: []string{"registry/package"},
			want: []PackageRef{
				{Registry: "registry", Name: "package", Version: ""},
			},
		},
		{
			name: "multiple args",
			args: []string{"registry1/package1", "registry2/package2@1.0.0"},
			want: []PackageRef{
				{Registry: "registry1", Name: "package1", Version: ""},
				{Registry: "registry2", Name: "package2", Version: "1.0.0"},
			},
		},
		{
			name: "empty args",
			args: []string{},
			want: []PackageRef{},
		},
		{
			name:    "invalid arg",
			args:    []string{"registry/package", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePackageArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePackageArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePackageArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultIncludePatterns(t *testing.T) {
	tests := []struct {
		name    string
		include []string
		want    []string
	}{
		{
			name:    "empty include",
			include: []string{},
			want:    []string{"*.yml", "*.yaml"},
		},
		{
			name:    "nil include",
			include: nil,
			want:    []string{"*.yml", "*.yaml"},
		},
		{
			name:    "existing patterns",
			include: []string{"*.md", "*.txt"},
			want:    []string{"*.md", "*.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDefaultIncludePatterns(tt.include)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDefaultIncludePatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}
