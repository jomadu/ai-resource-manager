package registry

import (
	"encoding/json"
	"fmt"
)

// Factory creates registries from config
type Factory interface {
	CreateRegistry(name string, config map[string]interface{}) (Registry, error)
}

// DefaultFactory is the default registry factory
type DefaultFactory struct{}

func (f *DefaultFactory) CreateRegistry(name string, config map[string]interface{}) (Registry, error) {
	regType, ok := config["type"].(string)
	if !ok {
		return nil, fmt.Errorf("registry type not specified")
	}

	switch regType {
	case "git":
		var gitConfig GitRegistryConfig
		if err := convertMapToStruct(config, &gitConfig); err != nil {
			return nil, err
		}
		return NewGitRegistry(name, gitConfig)
	default:
		return nil, fmt.Errorf("unsupported registry type: %s", regType)
	}
}

func convertMapToStruct(m map[string]interface{}, target interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
