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
	Install(ctx context.Context, dir, registry, ruleset, version string, files []types.File) error
	Uninstall(ctx context.Context, dir, registry, ruleset string) error
	ListInstalled(ctx context.Context, dir string) ([]Installation, error)
}
