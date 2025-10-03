# Sample Promptset Compilation

This directory contains a sample promptset (`sample-promptset.yml`) compiled to multiple build targets, demonstrating how the AI Rules Manager transforms a single promptset definition into platform-specific formats.

## Compilation Process

The compilation process takes a source promptset YAML file and generates platform-specific output files for each prompt. Unlike rulesets, promptsets compile to minimal, content-only files across all targets.

### Source Promptset Structure

The source `sample-promptset.yml` contains:
- **Promptset metadata**: ID, name, description
- **Prompts**: Each with optional properties like name, description, and body content
- **No enforcement levels**: Prompts don't have enforcement concepts like rules do
- **No scope**: Prompts apply broadly rather than to specific file patterns

### Key Compilation Principles

- **Content-only output**: All prompt files contain only the body content, no metadata or frontmatter
- **Currently identical across targets**: All build targets produce the same content (this may change as platform-specific features are added)
- **No platform-specific features yet**: Prompts don't currently have enforcement levels, file scopes, or other platform-specific metadata
- **Maximum simplicity**: Focus on pure content for easy integration with AI assistants

## Build Targets

### `/md/` - Markdown Format

**File Extension**: `.md`

**Purpose**: Simple markdown format containing only the prompt body content.

**Structure**:
- **Body**: Prompt content only

**Key Features**:
- Minimal, clean format
- No metadata or frontmatter
- Pure content focus
- Easy to read and use

**Example**:
```
This is how you review the code.
```

### `/cursor/` - Cursor IDE Format

**File Extension**: `.md`

**Purpose**: Simple format for Cursor IDE containing only the prompt body content.

**Structure**:
- **Body**: Prompt content only

**Key Features**:
- Minimal, clean format
- No metadata or frontmatter
- Pure content focus
- Easy to integrate with Cursor

**Example**:
```
This is how you review the code.
```

### `/amazonq/` - Amazon Q Format

**File Extension**: `.md`

**Purpose**: Simple format for Amazon Q AI assistant containing only the prompt body content.

**Structure**:
- **Body**: Prompt content only

**Key Features**:
- Minimal, clean format
- No metadata or frontmatter
- Pure content focus
- Easy to integrate with Amazon Q

**Example**:
```
This is how you review the code.
```

## Key Differences from Rulesets

### Simpler Compilation Process
- **No platform-specific metadata**: Unlike rulesets, prompts don't need different frontmatter for different platforms
- **No enforcement mapping**: No need to convert `must`/`should`/`may` to platform-specific properties
- **No file scope handling**: No need to convert scope arrays to different formats
- **No conditional logic**: All prompts compile identically regardless of their properties

### Content-Only Output
- **No metadata preservation**: Unlike rulesets, prompt metadata is not preserved in output files
- **No frontmatter blocks**: No YAML frontmatter in any build target
- **No headers or titles**: No markdown headers or enforcement indicators
- **Pure content focus**: Only the prompt body content is included

### Platform Consistency (Current State)
- **Currently identical across targets**: All build targets produce the same content
- **No target-specific features yet**: No `globs`, `applyTo`, `alwaysApply`, or other platform properties currently implemented
- **Universal format**: Same simple text format works for all AI assistants and tools
- **Future extensibility**: Platform-specific features may be added as the feature matures

## Usage

Each build target currently serves the same purpose with identical content:

- **Markdown**: Simple text files for documentation and general use
- **AmazonQ**: Direct integration with Amazon Q AI assistant
- **Cursor**: Direct integration with Cursor IDE
- **All targets**: Pure content for maximum compatibility with AI tools

The compilation process currently prioritizes simplicity and universal compatibility. As the promptset feature matures, platform-specific features may be added to provide more targeted integration with different AI assistants and development tools.
