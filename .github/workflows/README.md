# GitHub Actions Workflows

This directory contains GitHub Actions workflows for CI/CD automation.

## Workflows

### 1. `lint.yml` - Code Linting and Testing

**Triggers:**
- Push to `main` branch
- Pull requests targeting `main` branch

**Jobs:**

#### lint-backend
Lints the Go backend code:
- Runs `golangci-lint` for comprehensive Go linting
- Runs `go vet` for additional static analysis
- Checks code formatting with `gofmt`

#### lint-frontend
Lints the React frontend code:
- Runs ESLint on TypeScript/React code
- Performs TypeScript type checking

#### build-test
Validates that both applications build successfully:
- Builds the frontend with Vite
- Builds the Go backend binary

#### test-backend
Runs backend unit tests:
- Sets up PostgreSQL database
- Runs Go tests with race detection
- Generates and uploads code coverage reports

**Requirements:**
- Go 1.23+
- Node.js 20+
- PostgreSQL 16 (for tests)

---

### 2. `pr-checks.yml` - Pull Request Quality Checks

**Triggers:**
- Pull request events (opened, synchronize, reopened)

**Jobs:**

#### pr-title-check
Enforces conventional commit format for PR titles:
- Validates PR title follows semantic versioning format
- Allowed types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`
- Example: `feat: add user authentication`

#### code-quality
General code quality checks:
- Detects large files (>5MB)
- Scans for potential secrets using Gitleaks

#### dependency-review
Security check for dependencies:
- Reviews dependency changes
- Fails on moderate or higher severity vulnerabilities

#### changed-files
Determines which parts of the monorepo changed:
- Detects backend changes (`apps/backend/**`)
- Detects frontend changes (`apps/www/**`)

#### lint-backend-conditional
Runs backend linting only when backend files change:
- Uses `--new-from-rev` to lint only changed code
- More efficient for large codebases

#### lint-frontend-conditional
Runs frontend linting only when frontend files change:
- Lints only the changed files
- Faster feedback for frontend-only changes

---

## Setup Instructions

### For Repository Maintainers

1. **Enable GitHub Actions:**
   - Go to Settings → Actions → General
   - Enable "Allow all actions and reusable workflows"

2. **Configure Branch Protection:**
   - Go to Settings → Branches
   - Add rule for `main` branch:
     - ✅ Require a pull request before merging
     - ✅ Require status checks to pass before merging
     - Select required checks:
       - `Lint Go Backend`
       - `Lint React Frontend`
       - `Build Test`
       - `PR Title Check`
       - `Code Quality Check`

3. **Optional: Enable Codecov:**
   - Sign up at https://codecov.io
   - Add `CODECOV_TOKEN` to repository secrets
   - Coverage reports will be uploaded automatically

### For Contributors

1. **Before Submitting a PR:**
   ```bash
   # Backend
   cd apps/backend
   golangci-lint run
   go test ./...
   gofmt -s -w .

   # Frontend
   cd apps/www
   npm run lint
   npm run build
   ```

2. **PR Title Format:**
   Use conventional commit format:
   ```
   feat: add new feature
   fix: resolve bug in authentication
   docs: update README
   refactor: reorganize folder structure
   ```

3. **Common Issues:**
   - **Large files rejected:** Remove or use Git LFS
   - **Secrets detected:** Remove sensitive data, rotate credentials
   - **Linting failures:** Run linters locally first
   - **Type errors:** Fix TypeScript errors before pushing

---

## Extending the Workflows

### Adding New Jobs

Edit `lint.yml` or `pr-checks.yml`:

```yaml
new-job:
  name: New Job Name
  runs-on: ubuntu-latest
  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Your custom step
      run: echo "Hello World"
```

### Adding New Linters

**Backend:**
```yaml
- name: Run new linter
  working-directory: apps/backend
  run: your-linter-command
```

**Frontend:**
```yaml
- name: Run new linter
  working-directory: apps/www
  run: npm run your-lint-script
```

### Environment Variables

Add secrets in Settings → Secrets → Actions:
```yaml
env:
  SECRET_KEY: ${{ secrets.YOUR_SECRET }}
```

---

## Troubleshooting

### Workflow Fails on First Run
- Check Go version matches (1.23)
- Ensure Node.js version is 20+
- Verify `package-lock.json` exists

### Cache Issues
Clear cache in Actions tab → Caches → Delete

### Permission Errors
Ensure Actions have write permissions:
Settings → Actions → General → Workflow permissions → Read and write permissions

---

## Performance Optimization

Current optimizations:
- ✅ Dependency caching (Go modules, npm)
- ✅ Conditional job execution (only lint changed code)
- ✅ Parallel job execution
- ✅ Minimal container images

Future improvements:
- Add matrix testing for multiple Go/Node versions
- Implement custom caching strategies
- Add deployment workflows
