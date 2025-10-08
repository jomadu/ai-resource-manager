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
	InstallRuleset(ctx context.Context, registry, ruleset, version string, priority int, files []types.File) error
	InstallPromptset(ctx context.Context, registry, promptset, version string, files []types.File) error
	UninstallRuleset(ctx context.Context, registry, ruleset string) error
	UninstallPromptset(ctx context.Context, registry, promptset string) error
	ListInstalledRulesets(ctx context.Context) ([]Ruleset, error)
	ListInstalledPromptsets(ctx context.Context) ([]Promptset, error)
	IsRulesetInstalled(ctx context.Context, registry, ruleset string) (bool, string, error)
	IsPromptsetInstalled(ctx context.Context, registry, promptset string) (bool, string, error)
}
