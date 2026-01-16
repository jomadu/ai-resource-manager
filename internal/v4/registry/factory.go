package registry

import (
	"encoding/json"
	"fmt"
)

func CreateRegistry(name string, config map[string]interface{}) (Registry, error) {
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
