package storage

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "simple string",
			input:   "hello",
			wantErr: false,
		},
		{
			name:    "simple map",
			input:   map[string]interface{}{"url": "https://github.com/user/repo", "type": "git"},
			wantErr: false,
		},
		{
			name:    "complex map with arrays",
			input:   map[string]interface{}{"name": "clean-code", "includes": []string{"*.yml", "*.yaml"}, "excludes": []string{"test/**"}},
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: false,
		},
		{
			name:    "empty map",
			input:   map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "nested map",
			input:   map[string]interface{}{"config": map[string]interface{}{"nested": "value"}},
			wantErr: false,
		},
		{
			name:    "unmarshalable input",
			input:   make(chan int), // channels can't be marshaled to JSON
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GenerateKey(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GenerateKey() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("GenerateKey() unexpected error: %v", err)
				return
			}
			
			// Key should be non-empty hex string
			if key == "" {
				t.Errorf("GenerateKey() returned empty key")
			}
			
			// Key should be 64 characters (SHA256 hex)
			if len(key) != 64 {
				t.Errorf("GenerateKey() key length = %d, want 64", len(key))
			}
			
			// Key should be valid hex
			for _, c := range key {
				if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
					t.Errorf("GenerateKey() key contains invalid hex character: %c", c)
					break
				}
			}
		})
	}
}

func TestGenerateKey_Consistency(t *testing.T) {
	input := map[string]interface{}{
		"url":  "https://github.com/user/repo",
		"type": "git",
	}
	
	// Generate key multiple times
	key1, err1 := GenerateKey(input)
	key2, err2 := GenerateKey(input)
	key3, err3 := GenerateKey(input)
	
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("GenerateKey() unexpected errors: %v, %v, %v", err1, err2, err3)
	}
	
	// All keys should be identical
	if key1 != key2 || key2 != key3 {
		t.Errorf("GenerateKey() not consistent: %s, %s, %s", key1, key2, key3)
	}
}

func TestGenerateKey_DifferentInputs(t *testing.T) {
	input1 := map[string]interface{}{"url": "https://github.com/user/repo1", "type": "git"}
	input2 := map[string]interface{}{"url": "https://github.com/user/repo2", "type": "git"}
	
	key1, err1 := GenerateKey(input1)
	key2, err2 := GenerateKey(input2)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("GenerateKey() unexpected errors: %v, %v", err1, err2)
	}
	
	// Different inputs should produce different keys
	if key1 == key2 {
		t.Errorf("GenerateKey() same key for different inputs: %s", key1)
	}
}

func TestGenerateKey_OrderIndependent(t *testing.T) {
	// Maps with same content but potentially different iteration order
	input1 := map[string]interface{}{"a": "1", "b": "2", "c": "3"}
	input2 := map[string]interface{}{"c": "3", "a": "1", "b": "2"}
	
	key1, err1 := GenerateKey(input1)
	key2, err2 := GenerateKey(input2)
	
	if err1 != nil || err2 != nil {
		t.Fatalf("GenerateKey() unexpected errors: %v, %v", err1, err2)
	}
	
	// Same content should produce same key regardless of map order
	if key1 != key2 {
		t.Errorf("GenerateKey() different keys for same content: %s vs %s", key1, key2)
	}
}