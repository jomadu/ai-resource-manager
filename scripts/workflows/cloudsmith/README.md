# Cloudsmith Workflow Scripts

These scripts help you test ARM with Cloudsmith registry integration.

## Prerequisites

1. A Cloudsmith account with an API token
2. `CLOUDSMITH_OWNER` and `CLOUDSMITH_REPOSITORY` environment variables or `.env` file

## Setup

### Option 1: Use Environment Variables (Recommended)

Export `CLOUDSMITH_API_KEY` in your shell profile (e.g., `.zshenv`, `.bashrc`):

```bash
# In ~/.zshenv or ~/.bashrc
export CLOUDSMITH_API_KEY="your-cloudsmith-api-token"
export CLOUDSMITH_OWNER="your-username-or-org"
export CLOUDSMITH_REPOSITORY="your-repository"
```

Then reload your shell or run:
```bash
source ~/.zshenv  # or source ~/.bashrc
```

### Option 2: Use .env File

Copy the example file and fill in your values:

```bash
cd scripts/workflows/cloudsmith
cp .env.example .env
# Edit .env with your values
```

## Testing Environment Variable Expansion in .armrc

The `init-cloudsmith-sandbox.sh` script creates an `.armrc` file that uses environment variable expansion:

```ini
[registry https://api.cloudsmith.io/your-owner/your-repo]
token = ${CLOUDSMITH_API_KEY}
```

ARM's rcfile service (see `internal/rcfile/service.go`) automatically expands environment variables using the `${VAR_NAME}` syntax, so the token will be read from your environment at runtime.

## Scripts

### init-cloudsmith-sandbox.sh

Sets up a clean sandbox environment for testing Cloudsmith integration:

```bash
./init-cloudsmith-sandbox.sh
```

This script:
1. Builds the ARM binary
2. Creates a sandbox directory with a fresh ARM environment
3. Creates an `.armrc` file with `${CLOUDSMITH_API_KEY}` for env var expansion
4. Configures the Cloudsmith registry
5. Adds cursor and amazonq sinks

## Verifying Environment Variable Expansion

After running `init-cloudsmith-sandbox.sh`, you can verify that env vars are properly expanded:

1. Check the generated `.armrc`:
   ```bash
   cat sandbox/.armrc
   ```
   
   Should contain:
   ```ini
   [registry https://api.cloudsmith.io/your-owner/your-repo]
   token = ${CLOUDSMITH_API_KEY}
   ```

2. Test an installation:
   ```bash
   cd sandbox
   ./arm install ruleset cloudsmith-registry/ai-rules cursor-rules
   ```

3. If it works, environment variable expansion is working correctly!

## Troubleshooting

### Token not found

If you get authentication errors:

1. Verify `CLOUDSMITH_API_KEY` is exported:
   ```bash
   echo $CLOUDSMITH_API_KEY
   ```

2. Check that the registry section name in `.armrc` matches exactly:
   ```bash
   cat sandbox/arm.json  # Check registry URL
   cat sandbox/.armrc    # Check section name
   ```

3. Run with verbose mode:
   ```bash
   cd sandbox
   ./arm install ruleset cloudsmith-registry/ai-rules cursor-rules -v
   ```

### Invalid credentials

Verify your token works directly:

```bash
curl -H "X-Api-Key: $CLOUDSMITH_API_KEY" \
     "https://api.cloudsmith.io/v1/packages/$CLOUDSMITH_OWNER/$CLOUDSMITH_REPOSITORY/"
```

