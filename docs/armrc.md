# .armrc Configuration

The `.armrc` file provides authentication configuration for ARM registries that require tokens or credentials.

## File Location

ARM looks for `.armrc` in the following order:
1. Current working directory (`./.armrc`)
2. User home directory (`~/.armrc`)

## Format

The `.armrc` file uses INI format with registry-specific sections:

```ini
[registry https://hostname/project/id]
token = your_token_here

[registry https://hostname/group/id]
token = your_token_here
```

**Note**: The full URL including protocol is required to avoid ambiguity between HTTP/HTTPS registries and custom ports.

## GitLab Registry Authentication

GitLab registries require explicit token authentication because they use GitLab's Package Registry API.

### Project-Level Packages

For project-level packages, use the project ID:

```ini
[registry https://gitlab.example.com/project/123]
token = glpat-xxxxxxxxxxxxxxxxxxxx
```

### Group-Level Packages

For group-level packages, use the group ID:

```ini
[registry https://gitlab.example.com/group/456]
token = glpat-xxxxxxxxxxxxxxxxxxxx
```

### Token Types

GitLab supports several token types:

- **Personal Access Token**: `glpat-xxxxxxxxxxxxxxxxxxxx`
- **Project Access Token**: `glpat-xxxxxxxxxxxxxxxxxxxx`
- **Group Access Token**: `glpat-xxxxxxxxxxxxxxxxxxxx`
- **Deploy Token**: `gldt-xxxxxxxxxxxxxxxxxxxx`

Required scopes: `read_api`, `read_repository`

### Example Configuration

Complete `.armrc` example:

```ini
# Production GitLab instance
[registry https://gitlab.company.com/project/5950]
token = glpat-abc123def456ghi789

# GitLab.com project
[registry https://gitlab.example.com/project/12345]
token = glpat-xyz789uvw456rst123

# Group-level packages
[registry https://gitlab.company.com/group/100]
token = glpat-group-token-here

# Internal HTTP registry with custom port
[registry http://internal-gitlab.company.com:8080/project/999]
token = glpat-internal-token-here

# Cloudsmith registry
[registry https://api.cloudsmith.io/myorg/ai-rules]
token = ckcy_abc123def456ghi789jkl012mno345pqr678stu901vwx234yz
```

## Security Notes

- **File Permissions**: Ensure `.armrc` has restricted permissions (`chmod 600 .armrc`)
- **Environment Variables**: Use environment variable substitution for CI/CD:
  ```ini
  [registry https://gitlab.example.com/project/123]
  token = ${GITLAB_TOKEN}
  ```
- **Gitignore**: Add `.armrc` to your `.gitignore` to avoid committing tokens
- **Local vs Global**: Use local `.armrc` for project-specific tokens, global `~/.armrc` for personal tokens

## Registry Identification

The registry section name must match the full URL and path used in the registry configuration:

- Registry URL: `https://gitlab.example.com` with `--project-id 123`
- Section name: `[registry https://gitlab.example.com/project/123]`

- Registry URL: `https://gitlab.example.com` with `--group-id 456`
- Section name: `[registry https://gitlab.example.com/group/456]`

- Registry URL: `http://internal-gitlab.company.com:8080` with `--project-id 999`
- Section name: `[registry http://internal-gitlab.company.com:8080/project/999]`

## Troubleshooting

### Common Issues

1. **Section not found**: Ensure the section name exactly matches the registry configuration
2. **Invalid token**: Verify token has correct scopes and hasn't expired
3. **Permission denied**: Check file permissions on `.armrc`
4. **Environment variable not expanded**: Ensure the variable is set when ARM runs

### Debug Authentication

Use verbose mode to see authentication details:

```bash
arm install my-gitlab/ai-rules --sinks cursor -v
```

### Test Token

Verify your token works with GitLab API:

```bash
curl --header "PRIVATE-TOKEN: your_token" \
     "https://gitlab.example.com/api/v4/projects/123/packages"
```
