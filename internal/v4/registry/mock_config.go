package registry

import "context"

type mockConfigManager struct {
	sections map[string]map[string]string
}

func newMockConfigManager() *mockConfigManager {
	return &mockConfigManager{
		sections: make(map[string]map[string]string),
	}
}

func (m *mockConfigManager) SetValue(section, key, value string) {
	if m.sections[section] == nil {
		m.sections[section] = make(map[string]string)
	}
	m.sections[section][key] = value
}

func (m *mockConfigManager) GetAllSections(ctx context.Context) (map[string]map[string]string, error) {
	return m.sections, nil
}

func (m *mockConfigManager) GetSection(ctx context.Context, section string) (map[string]string, error) {
	if sec, ok := m.sections[section]; ok {
		return sec, nil
	}
	return nil, nil
}

func (m *mockConfigManager) GetValue(ctx context.Context, section, key string) (string, error) {
	if sec, ok := m.sections[section]; ok {
		if val, ok := sec[key]; ok {
			return val, nil
		}
	}
	return "", nil
}
