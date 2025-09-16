package installer

import (
	"context"
	"os"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func isEmpty(dir string) bool {
	entries, err := os.ReadDir(dir)
	return err == nil && len(entries) == 0
}

// Installer manages physical file deployment to sink directories.
type Installer interface {
	Install(ctx context.Context, registry, ruleset, version string, priority int, files []types.File) error
	Uninstall(ctx context.Context, registry, ruleset string) error
	ListInstalled(ctx context.Context) ([]Ruleset, error)
	IsInstalled(ctx context.Context, registry, ruleset string) (bool, string, error)
}
