package registry

// MockFactory is a test factory that returns a mock registry
type MockFactory struct {
	registry Registry
}

func NewMockFactory(reg Registry) *MockFactory {
	return &MockFactory{registry: reg}
}

func (f *MockFactory) CreateRegistry(name string, config map[string]interface{}) (Registry, error) {
	return f.registry, nil
}
