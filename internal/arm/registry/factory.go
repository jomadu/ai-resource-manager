package registry

import (
	"encoding/json"
	"fmt"

	"github.com/jomadu/ai-resource-manager/internal/arm/config"
)

// Factory creates registries from config
type Factory interface {
	CreateRegistry(name string, config map[string]interface{}) (Registry, error)
}

// DefaultFactory is the default registry factory
type DefaultFactory struct{}

func (f *DefaultFactory) CreateRegistry(name string, cfg map[string]interface{}) (Registry, error) {
	regType, ok := cfg["type"].(string)
	if !ok {
		return nil, fmt.Errorf("registry type not specified")
	}

	switch regType {
	case "git":
		var gitConfig GitRegistryConfig
		if err := convertMapToStruct(cfg, &gitConfig); err != nil {
			return nil, err
		}
		return NewGitRegistry(name, gitConfig)
	case "gitlab":
		var gitlabConfig GitLabRegistryConfig
		if err := convertMapToStruct(cfg, &gitlabConfig); err != nil {
			return nil, err
		}
		configMgr := config.NewFileManager()
		return NewGitLabRegistry(name, gitlabConfig, configMgr)
	case "cloudsmith":
		var cloudsmithConfig CloudsmithRegistryConfig
		if err := convertMapToStruct(cfg, &cloudsmithConfig); err != nil {
			return nil, err
		}
		configMgr := config.NewFileManager()
		return NewCloudsmithRegistry(name, cloudsmithConfig, configMgr)
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
