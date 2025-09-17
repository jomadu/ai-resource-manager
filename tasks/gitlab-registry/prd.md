# GitLab Registry Support - Product Requirements Document

## Overview

Add support for GitLab Generic Package Registry as a registry type in ARM, enabling users to store and distribute AI rulesets through GitLab's package management infrastructure.

## Problem Statement

Currently, ARM only supports Git-based registries, which have significant limitations:

### Path Fragility in Git Registries
Git registries require brittle `--include`/`--exclude` path targeting (e.g., `--include "rules/cursor/*.mdc"`). These paths are fragile because:
- Repository maintainers can reorganize files, breaking existing installations
- Consumers are tightly coupled to internal repository structure
- No semantic way to request rulesets - users must know exact directory layouts
- Path changes under consumers' feet cause silent failures or missing updates

### GitLab Package Registry Benefits
Organizations using GitLab want to leverage their existing package registry infrastructure for AI ruleset distribution, benefiting from:

- **Semantic targeting**: Install by package name (`cursor-rules`) not file paths
- **Encapsulated structure**: Package internals hidden from consumers
- **Reorganization resilience**: Repository changes don't break installations
- Centralized package management alongside other artifacts
- Built-in access control and authentication
- Package versioning and metadata
- Integration with existing CI/CD pipelines
- Bandwidth optimization through GitLab's CDN

## Goals

### Primary Goals
- Support GitLab Generic Package Registry as a registry type
- Enable authentication via multiple token types (Personal, Project, Deploy, CI/CD)
- Provide seamless installation and update workflows
- Maintain compatibility with existing ARM features (URF, versioning, etc.)

### Secondary Goals
- Support both project-level and group-level registries
- Enable configuration of custom GitLab instances (self-hosted)
- Provide clear error messages for authentication and access issues

## User Stories

### As a DevOps Engineer
- I want to publish AI rulesets to our GitLab package registry so they're managed alongside our other artifacts
- I want to use deploy tokens for automated CI/CD access to rulesets
- I want to configure ARM to use our self-hosted GitLab instance

### As a Developer
- I want to install AI rulesets from GitLab registries using familiar ARM commands
- I want ARM to handle GitLab authentication transparently
- I want to receive clear error messages when authentication fails

### As a Team Lead
- I want to control access to AI rulesets using GitLab's existing permission system
- I want to track ruleset usage through GitLab's package registry analytics

## Requirements

### Functional Requirements

#### Registry Configuration
- Support `gitlab` registry type in `arm config registry add`
- Accept GitLab project/group URLs or IDs
- Configure API version (default to v4)
- Support custom GitLab instances (self-hosted)

#### Authentication
- Store tokens in project rc file (`.armrc`) excluded from version control
- Support environment variable expansion in `.armrc` (e.g., `token=${GITLAB_TOKEN}`)
- Support any GitLab token type that provides package registry access

#### Package Operations
- List available rulesets from GitLab registry
- Download ruleset packages with version resolution
- Support semantic versioning constraints
- Handle GitLab's package metadata format
- Maintain compatibility with existing ARM workflows

#### Error Handling
- Authentication error messages
- Network connectivity error handling
- Version resolution error reporting

### Non-Functional Requirements

#### Security
- Never store tokens in shared project configuration files (arm.json)
- Store tokens in project rc file (`.armrc`)
- Use HTTPS for all GitLab API communications

#### Performance
- Integrate with existing ARM caching mechanisms
- Respect GitLab API rate limits

#### Usability
- Consistent CLI interface with existing registry types
- Standard error handling and reporting

## Success Metrics

- Users can successfully configure GitLab registries
- Authentication works with GitLab tokens
- Package installation performance matches Git registry performance
- Secure token handling via .armrc file

## Out of Scope

- GitLab Container Registry support
- GitLab Maven/npm/other package types
- GitLab CI/CD pipeline integration
- Automatic token generation/management
- GitLab webhook integration for auto-updates

## Dependencies

- GitLab API v4 compatibility
- HTTP client with proper SSL/TLS support
- Semantic versioning library compatibility
- File system permissions for secure rc file storage

## Risks and Mitigations

### Risk: Token Security
**Mitigation**: Use project rc file, ensure .armrc is excluded from version control

### Risk: GitLab API Changes
**Mitigation**: Use versioned API endpoints, implement error handling
