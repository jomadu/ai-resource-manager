# Unit Test Checklist

## High Priority

- [x] `resolver.Constraint` - Version constraint with semantic versioning
- [x] `lockfile.LockFile` - Complete arm.lock structure
- [x] `lockfile.Entry` - Single ruleset entry in lock file
- [ ] `manifest.Manifest` - Complete arm.json structure
- [ ] `manifest.Entry` - Single ruleset entry in manifest
- [ ] `config.RCConfig` - Complete .armrc.json structure

## Medium Priority

- [ ] `arm.VersionRef` - Version references (tags, branches, commits)
- [ ] `arm.ContentSelector` - Include/exclude patterns for filtering
- [ ] `arm.OutdatedRuleset` - Outdated ruleset information
