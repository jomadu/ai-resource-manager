# Changelog

All notable changes to the Clean Code ruleset will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.0] - 2024-01-15

### Added
- New rule: `writeComprehensiveTests` - Comprehensive testing requirements with 80%+ coverage target
- Support for additional programming languages: Rust (`**/*.rs`), C++ (`**/*.cpp`), C (`**/*.c`)
- Enhanced rule descriptions with more detailed examples and guidance
- Breaking change documentation within rule bodies

### Changed
- **BREAKING**: Renamed `meaningfulNames` to `useMeaningfulNames`
- **BREAKING**: Renamed `smallFunctions` to `keepFunctionsSmall`
- **BREAKING**: Renamed `avoidComments` to `writeSelfDocumentingCode`
- **BREAKING**: Renamed `singleResponsibility` to `singleResponsibilityPrinciple`
- **BREAKING**: Renamed `avoidDuplication` to `eliminateDuplication`
- **BREAKING**: Renamed `consistentFormatting` to `consistentCodeStyle`
- **BREAKING**: Renamed `errorHandling` to `robustErrorHandling`
- **BREAKING**: Upgraded `writeSelfDocumentingCode` enforcement from SHOULD to MUST
- **BREAKING**: Upgraded `singleResponsibilityPrinciple` enforcement from SHOULD to MUST
- **BREAKING**: Upgraded `eliminateDuplication` enforcement from MAY to SHOULD
- **BREAKING**: Upgraded `consistentCodeStyle` enforcement from SHOULD to MUST
- **BREAKING**: Reduced function length limit from 20 to 15 lines in `keepFunctionsSmall`
- **BREAKING**: Increased priority of `keepFunctionsSmall` from 90 to 95
- **BREAKING**: Increased priority of `robustErrorHandling` from 85 to 90
- Updated ruleset description to "Essential clean code principles for maintainable software - Major refactor with enhanced rules"

### Migration Guide
Users upgrading from v1.x will need to:
1. Update any references to old rule IDs
2. Review enforcement level changes and adjust compliance
3. Update tooling to handle new language support
4. Implement new testing requirements

## [1.1.0] - 2024-01-10

### Added
- New rule: `consistentFormatting` - Code formatting and style consistency requirements
- New rule: `errorHandling` - Proper error handling practices and guidelines

### Changed
- No breaking changes to existing rules
- All existing rules remain unchanged

## [1.0.1] - 2024-01-05

### Fixed
- Minor documentation improvements
- Bug fixes and corrections

### Changed
- Identical to v1.0.0
- No functional changes

## [1.0.0] - 2024-01-01

### Added
- Initial release of Clean Code ruleset
- Core rule: `meaningfulNames` - Use intention-revealing names (MUST, Priority 100)
- Core rule: `smallFunctions` - Keep functions small and focused (MUST, Priority 90)
- Core rule: `avoidComments` - Write self-documenting code (SHOULD, Priority 80)
- Core rule: `singleResponsibility` - Single responsibility principle (SHOULD, Priority 85)
- Core rule: `avoidDuplication` - Don't repeat yourself (DRY) (MAY, Priority 75)
- Support for Python, JavaScript, TypeScript, Java, and Go
- Comprehensive rule descriptions with examples and best practices

### Changed
- N/A (initial release)

### Fixed
- N/A (initial release)

### Removed
- N/A (initial release)

---

## Version Compatibility Matrix

| From Version | To Version | Compatibility | Migration Required |
|--------------|------------|---------------|-------------------|
| 1.0.0        | 1.0.1      | ✅ Full       | No                |
| 1.0.x        | 1.1.0      | ✅ Full       | No                |
| 1.x.x        | 2.0.0      | ❌ Breaking   | Yes               |

## Rule Evolution

| Rule Name | v1.0.0 | v1.1.0 | v2.0.0 | Notes |
|-----------|--------|--------|--------|-------|
| Meaningful Names | ✅ | ✅ | ✅ | Renamed in v2.0.0 |
| Small Functions | ✅ | ✅ | ✅ | Renamed, stricter limits in v2.0.0 |
| Avoid Comments | ✅ | ✅ | ✅ | Renamed, enforcement upgraded in v2.0.0 |
| Single Responsibility | ✅ | ✅ | ✅ | Renamed, enforcement upgraded in v2.0.0 |
| Avoid Duplication | ✅ | ✅ | ✅ | Renamed, enforcement upgraded in v2.0.0 |
| Consistent Formatting | ❌ | ✅ | ✅ | New in v1.1.0, renamed in v2.0.0 |
| Error Handling | ❌ | ✅ | ✅ | New in v1.1.0, renamed in v2.0.0 |
| Comprehensive Tests | ❌ | ❌ | ✅ | New in v2.0.0 |

## Contributing

When contributing to this ruleset:

- **Patch versions (1.0.x)**: Bug fixes and minor corrections only
- **Minor versions (1.x.0)**: New rules and features, backward compatible
- **Major versions (x.0.0)**: Breaking changes, rule renames, enforcement changes

Always document breaking changes clearly and provide migration guidance.
