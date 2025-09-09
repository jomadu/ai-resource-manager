package cache

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		wantErr bool
	}{
		{
			name: "simple string",
			obj:  "test",
		},
		{
			name: "struct",
			obj:  struct{ Name string }{Name: "test"},
		},
		{
			name: "map",
			obj:  map[string]string{"key": "value"},
		},
		{
			name: "nil",
			obj:  nil,
		},
		{
			name:    "unmarshalable",
			obj:     make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateKey(tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != 64 {
				t.Errorf("GenerateKey() = %v, want 64 character hash", got)
			}
		})
	}
}

func TestGenerateKeyDeterministic(t *testing.T) {
	obj := map[string]string{"key": "value"}
	key1, err := GenerateKey(obj)
	if err != nil {
		t.Fatal(err)
	}
	key2, err := GenerateKey(obj)
	if err != nil {
		t.Fatal(err)
	}
	if key1 != key2 {
		t.Errorf("GenerateKey() not deterministic: %v != %v", key1, key2)
	}
}
