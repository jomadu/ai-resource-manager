# Changelog

All notable changes to the Code Review promptset will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.0] - 2024-01-15

### Added
- New prompt: `devopsAndDeploymentReview` - DevOps practices and operational excellence assessment
- Enhanced security audit with threat modeling and compliance focus
- Advanced performance analysis with profiling and scalability assessment
- Comprehensive code quality assessment with architecture focus
- Inclusive design review with WCAG 2.1 AA compliance
- Detailed prompt descriptions with structured guidance

### Changed
- **BREAKING**: Renamed `securityReview` to `comprehensiveSecurityAudit`
- **BREAKING**: Renamed `performanceReview` to `performanceOptimizationAnalysis`
- **BREAKING**: Renamed `codeQualityReview` to `codeQualityExcellence`
- **BREAKING**: Renamed `accessibilityReview` to `inclusiveDesignReview`
- **BREAKING**: Updated promptset name from "Code Review" to "Advanced Code Review"
- **BREAKING**: Enhanced all existing prompts with significantly more detailed content
- **BREAKING**: Restructured prompt bodies with organized sections and categories
- Updated promptset description to "Comprehensive code review prompts for quality assurance - Major refactor with enhanced prompts"

### Migration Guide
Users upgrading from v1.x will need to:
1. Update any references to old prompt IDs
2. Review new prompt content structure and adapt tooling
3. Update AI assistant configurations to handle enhanced prompts
4. Consider the new DevOps and deployment review capabilities

## [1.1.0] - 2024-01-10

### Added
- New prompt: `codeQualityReview` - Code quality, maintainability, and best practices evaluation
- New prompt: `accessibilityReview` - Accessibility compliance and inclusive design assessment

### Changed
- No breaking changes to existing prompts
- All existing prompts remain unchanged

## [1.0.1] - 2024-01-05

### Fixed
- Minor documentation improvements
- Bug fixes and corrections

### Changed
- Identical to v1.0.0
- No functional changes

## [1.0.0] - 2024-01-01

### Added
- Initial release of Code Review promptset
- Core prompt: `securityReview` - Security vulnerabilities and best practices focus
- Core prompt: `performanceReview` - Performance bottlenecks and optimization analysis
- Comprehensive prompt descriptions with detailed guidance
- Focus on security and performance review aspects

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

## Prompt Evolution

| Prompt Name | v1.0.0 | v1.1.0 | v2.0.0 | Notes |
|-------------|--------|--------|--------|-------|
| Security Review | ✅ | ✅ | ✅ | Renamed to "Comprehensive Security Audit" in v2.0.0 |
| Performance Review | ✅ | ✅ | ✅ | Renamed to "Performance Optimization Analysis" in v2.0.0 |
| Code Quality Review | ❌ | ✅ | ✅ | New in v1.1.0, renamed to "Code Quality Excellence" in v2.0.0 |
| Accessibility Review | ❌ | ✅ | ✅ | New in v1.1.0, renamed to "Inclusive Design Review" in v2.0.0 |
| DevOps and Deployment Review | ❌ | ❌ | ✅ | New in v2.0.0 |

## Prompt Content Evolution

### v1.0.0 - Basic Prompts
- Simple bullet-point lists
- Focus on security and performance
- Basic guidance structure

### v1.1.0 - Enhanced Coverage
- Added code quality and accessibility
- Maintained simple structure
- Expanded review coverage

### v2.0.0 - Comprehensive Analysis
- **Structured sections** with clear categories
- **Detailed subcategories** for each review area
- **Professional formatting** with markdown headers
- **Comprehensive coverage** including DevOps and deployment
- **Enhanced descriptions** with specific focus areas

## Usage Recommendations

### For New Projects
- **Recommended**: Use the latest version (2.0.0) for comprehensive code review
- **Benefits**: Most detailed prompts with structured guidance and comprehensive coverage

### For Existing Projects
- **v1.x.x**: Continue using if simple prompts are sufficient
- **v2.0.0**: Upgrade when ready for comprehensive, structured review prompts

### For AI Assistants
- **v1.x.x**: Simpler integration, basic prompt structure
- **v2.0.0**: More sophisticated integration, structured prompt handling

## Contributing

When contributing to this promptset:

- **Patch versions (1.0.x)**: Bug fixes and minor corrections only
- **Minor versions (1.x.0)**: New prompts and features, backward compatible
- **Major versions (x.0.0)**: Breaking changes, prompt renames, content restructuring

Always document breaking changes clearly and provide migration guidance.
