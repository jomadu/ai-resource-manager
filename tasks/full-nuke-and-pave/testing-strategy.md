# ARM Testing Strategy

## Current Test Coverage Analysis

### Existing Tests
- `cmd/arm/cache_test.go` - Duration parsing utility tests ✅
- `internal/types/version_test.go` - Version type tests ✅
- `internal/config/manager_test.go` - Config management tests ✅
- `internal/manifest/manager_test.go` - Manifest management tests ✅
- `internal/lockfile/checksum_test.go` - Checksum validation tests ✅
- `internal/lockfile/manager_test.go` - Lockfile management tests ✅
- `internal/registry/git_registry_test.go` - Git registry tests ✅
- `internal/arm/sink_sync_test.go` - Sink synchronization tests ✅
- `internal/arm/integration_test.go` - Integration tests ✅
- `internal/installer/installer_test.go` - File installation tests ✅

### Coverage Gaps
- **Service Layer:** No unit tests for main `ArmService` methods
- **CLI Layer:** No tests for command parsing and output formatting
- **Error Scenarios:** Limited error path testing
- **Edge Cases:** Missing boundary condition tests

## Testing Pyramid Strategy

### Unit Tests (70% of test effort)
**Goal:** Test individual components in isolation with mocked dependencies

#### Service Layer Tests
```go
// internal/arm/service_test.go
func TestArmService_InstallRuleset(t *testing.T) {
    tests := []struct {
        name    string
        request InstallRequest
        setup   func(*mocks.MockDependencies)
        wantErr bool
    }{
        {
            name: "successful install",
            request: InstallRequest{
                Registry: "test-registry",
                Ruleset:  "test-ruleset",
                Version:  "1.0.0",
            },
            setup: func(m *mocks.MockDependencies) {
                m.ManifestManager.EXPECT().GetRawRegistries(gomock.Any()).Return(validRegistries, nil)
                m.RegistryFactory.EXPECT().NewRegistry(gomock.Any(), gomock.Any()).Return(mockRegistry, nil)
                // ... more expectations
            },
            wantErr: false,
        },
        {
            name: "registry not found",
            request: InstallRequest{
                Registry: "nonexistent",
                Ruleset:  "test-ruleset",
            },
            setup: func(m *mocks.MockDependencies) {
                m.ManifestManager.EXPECT().GetRawRegistries(gomock.Any()).Return(emptyRegistries, nil)
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            deps := setupMockDependencies(t)
            tt.setup(deps)

            service := NewArmService(ServiceDependencies{
                ConfigManager:   deps.ConfigManager,
                ManifestManager: deps.ManifestManager,
                LockfileManager: deps.LockfileManager,
                RegistryFactory: deps.RegistryFactory,
            })

            err := service.InstallRuleset(context.Background(), tt.request)
            if (err != nil) != tt.wantErr {
                t.Errorf("InstallRuleset() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

#### Version Service Tests
```go
// internal/arm/version_service_test.go
func TestVersionService_ResolveVersion(t *testing.T) {
    tests := []struct {
        name     string
        request  InstallRequest
        mockResp *types.ResolvedVersion
        mockErr  error
        want     *ResolvedVersion
        wantErr  bool
    }{
        {
            name: "resolve latest version",
            request: InstallRequest{
                Registry: "test-registry",
                Ruleset:  "test-ruleset",
                Version:  "latest",
            },
            mockResp: &types.ResolvedVersion{
                Version: types.Version{Version: "1.2.3", Display: "v1.2.3"},
            },
            want: &ResolvedVersion{
                Version: types.Version{Version: "1.2.3", Display: "v1.2.3"},
            },
            wantErr: false,
        },
        {
            name: "expand version shorthand",
            request: InstallRequest{
                Registry: "test-registry",
                Ruleset:  "test-ruleset",
                Version:  "1.2",
            },
            // Should expand "1.2" to "^1.2.0"
            mockResp: &types.ResolvedVersion{
                Version: types.Version{Version: "1.2.5", Display: "v1.2.5"},
            },
            want: &ResolvedVersion{
                Version: types.Version{Version: "1.2.5", Display: "v1.2.5"},
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRegistry := mocks.NewMockRegistry(t)
            mockFactory := mocks.NewMockRegistryFactory(t)

            mockFactory.EXPECT().GetRegistry(tt.request.Registry).Return(mockRegistry, nil)
            mockRegistry.EXPECT().ResolveVersion(gomock.Any(), gomock.Any()).Return(tt.mockResp, tt.mockErr)

            service := NewVersionService(mockFactory)
            got, err := service.ResolveVersion(context.Background(), tt.request)

            if (err != nil) != tt.wantErr {
                t.Errorf("ResolveVersion() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ResolveVersion() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

#### CLI Parser Tests
```go
// cmd/arm/parser_test.go
func TestParseRulesetArg(t *testing.T) {
    tests := []struct {
        name    string
        arg     string
        want    RulesetRef
        wantErr bool
    }{
        {
            name: "registry/ruleset format",
            arg:  "ai-rules/amazonq-rules",
            want: RulesetRef{
                Registry: "ai-rules",
                Name:     "amazonq-rules",
                Version:  "",
            },
            wantErr: false,
        },
        {
            name: "registry/ruleset@version format",
            arg:  "ai-rules/amazonq-rules@1.2.3",
            want: RulesetRef{
                Registry: "ai-rules",
                Name:     "amazonq-rules",
                Version:  "1.2.3",
            },
            wantErr: false,
        },
        {
            name:    "invalid format - no registry",
            arg:     "just-ruleset",
            wantErr: true,
        },
        {
            name:    "empty argument",
            arg:     "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseRulesetArg(tt.arg)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseRulesetArg() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseRulesetArg() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests (20% of test effort)
**Goal:** Test component interactions and end-to-end workflows

#### End-to-End Workflow Tests
```go
// internal/arm/integration_test.go (enhance existing)
func TestCompleteWorkflow(t *testing.T) {
    // Setup test environment
    tempDir := t.TempDir()
    setupTestRegistry(t, tempDir)

    service := setupTestService(t, tempDir)

    // Test complete workflow: install -> list -> update -> uninstall
    t.Run("install ruleset", func(t *testing.T) {
        err := service.InstallRuleset(context.Background(), InstallRequest{
            Registry: "test-registry",
            Ruleset:  "test-ruleset",
            Version:  "1.0.0",
            Include:  []string{"**/*.md"},
        })
        require.NoError(t, err)
    })

    t.Run("list installed rulesets", func(t *testing.T) {
        installed, err := service.List(context.Background())
        require.NoError(t, err)
        require.Len(t, installed, 1)
        assert.Equal(t, "test-registry", installed[0].Registry)
        assert.Equal(t, "test-ruleset", installed[0].Name)
    })

    t.Run("update ruleset", func(t *testing.T) {
        err := service.UpdateRuleset(context.Background(), "test-registry", "test-ruleset")
        require.NoError(t, err)
    })

    t.Run("uninstall ruleset", func(t *testing.T) {
        err := service.Uninstall(context.Background(), "test-registry", "test-ruleset")
        require.NoError(t, err)

        // Verify removal
        installed, err := service.List(context.Background())
        require.NoError(t, err)
        assert.Len(t, installed, 0)
    })
}
```

#### Error Scenario Tests
```go
func TestErrorScenarios(t *testing.T) {
    service := setupTestService(t, t.TempDir())

    t.Run("install nonexistent registry", func(t *testing.T) {
        err := service.InstallRuleset(context.Background(), InstallRequest{
            Registry: "nonexistent",
            Ruleset:  "test-ruleset",
        })
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "registry nonexistent not configured")
    })

    t.Run("install with network failure", func(t *testing.T) {
        // Test with mock registry that returns network error
        // Verify graceful error handling
    })

    t.Run("install with corrupted lockfile", func(t *testing.T) {
        // Test recovery from corrupted state
    })
}
```

### Contract Tests (10% of test effort)
**Goal:** Test CLI interface contracts and output formats

#### CLI Command Tests
```go
// cmd/arm/commands_test.go
func TestInstallCommand(t *testing.T) {
    tests := []struct {
        name     string
        args     []string
        flags    map[string]string
        mockResp error
        wantErr  bool
        wantOut  string
    }{
        {
            name: "install single ruleset",
            args: []string{"ai-rules/test-ruleset@1.0.0"},
            flags: map[string]string{
                "include": "**/*.md",
            },
            mockResp: nil,
            wantErr:  false,
        },
        {
            name:     "install with invalid format",
            args:     []string{"invalid-format"},
            wantErr:  true,
            wantOut:  "invalid ruleset format",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockService := mocks.NewMockService(t)
            if tt.mockResp != nil {
                mockService.EXPECT().InstallRuleset(gomock.Any(), gomock.Any()).Return(tt.mockResp)
            }

            cmd := newInstallCommand(mockService)

            // Set flags
            for flag, value := range tt.flags {
                cmd.Flags().Set(flag, value)
            }

            // Capture output
            var buf bytes.Buffer
            cmd.SetOut(&buf)
            cmd.SetErr(&buf)

            err := cmd.RunE(cmd, tt.args)

            if (err != nil) != tt.wantErr {
                t.Errorf("command error = %v, wantErr %v", err, tt.wantErr)
            }

            if tt.wantOut != "" {
                output := buf.String()
                assert.Contains(t, output, tt.wantOut)
            }
        })
    }
}
```

## Mock Generation Strategy

### Use gomock for Interface Mocking
```bash
# Generate mocks for all service interfaces
go generate ./...

//go:generate mockgen -source=service.go -destination=mocks/mock_service.go
//go:generate mockgen -source=../config/manager.go -destination=mocks/mock_config_manager.go
//go:generate mockgen -source=../manifest/manager.go -destination=mocks/mock_manifest_manager.go
```

### Test Utilities
```go
// internal/arm/testutil/setup.go
func SetupMockDependencies(t *testing.T) *MockDependencies {
    ctrl := gomock.NewController(t)

    return &MockDependencies{
        ConfigManager:   mocks.NewMockConfigManager(ctrl),
        ManifestManager: mocks.NewMockManifestManager(ctrl),
        LockfileManager: mocks.NewMockLockfileManager(ctrl),
        RegistryFactory: mocks.NewMockRegistryFactory(ctrl),
    }
}

func SetupTestService(t *testing.T, tempDir string) *ArmService {
    // Setup real file-based managers for integration tests
    return NewArmService(ServiceDependencies{
        ConfigManager:   config.NewFileManager(filepath.Join(tempDir, ".armrc.json")),
        ManifestManager: manifest.NewFileManager(filepath.Join(tempDir, "arm.json")),
        LockfileManager: lockfile.NewFileManager(filepath.Join(tempDir, "arm-lock.json")),
        RegistryFactory: registry.NewFactory(),
    })
}
```

## Test Coverage Goals

### Minimum Coverage Targets
- **Service Layer:** 85% line coverage
- **CLI Layer:** 70% line coverage
- **Overall Project:** 75% line coverage

### Coverage Measurement
```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Check coverage percentage
go tool cover -func=coverage.out | grep total
```

## Continuous Testing Strategy

### Pre-commit Hooks
```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: go-test
        name: go test
        entry: go test ./...
        language: system
        pass_filenames: false

      - id: go-test-coverage
        name: go test coverage
        entry: bash -c 'go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep total | awk "{if(\$3+0 < 75) exit 1}"'
        language: system
        pass_filenames: false
```

### CI Pipeline Tests
```yaml
# .github/workflows/test.yml
name: Test
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'

      - name: Run tests
        run: go test -v -coverprofile=coverage.out ./...

      - name: Check coverage
        run: |
          go tool cover -func=coverage.out | grep total
          go tool cover -func=coverage.out | grep total | awk '{if($3+0 < 75) exit 1}'

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

## Implementation Timeline

### Week 1: Foundation
- Set up mock generation
- Create test utilities and helpers
- Write service layer unit tests

### Week 2: CLI & Integration
- Add CLI command tests
- Enhance integration tests
- Add error scenario coverage

### Week 3: Coverage & Quality
- Achieve coverage targets
- Add performance benchmarks
- Set up CI pipeline integration

## Success Metrics

- [ ] 85%+ service layer test coverage
- [ ] 70%+ CLI layer test coverage
- [ ] 75%+ overall project coverage
- [ ] All critical workflows have integration tests
- [ ] All error paths have test coverage
- [ ] CI pipeline enforces coverage requirements
