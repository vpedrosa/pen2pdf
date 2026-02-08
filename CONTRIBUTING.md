# Contributing Guide - Sinttik

## Conventional Commits

This project follows the [Conventional Commits](https://www.conventionalcommits.org/en/) specification for commit messages. This allows automatic changelog generation and makes the project history easier to understand.

### Message Format

```
<type> #<issue-number>: <description>

[optional body]

[optional footer(s)]
```

**Example:** `feat #45: add email monitoring capability`

> **IMPORTANT:** Every commit MUST reference a GitLab issue with `#number` and contain ONLY the changes related to that issue. Do not mix changes from multiple issues in a single commit.

### Commit Types

| Type | Description | Example |
|------|-------------|---------|
| `feat` | New feature | `feat #45: add task scheduler` |
| `fix` | Bug fix | `fix #78: resolve null pointer in handler` |
| `docs` | Documentation changes | `docs #92: update installation steps` |
| `style` | Formatting, whitespace, semicolons (no code changes) | `style #15: format code with gofmt` |
| `refactor` | Refactoring (no functional change or fix) | `refactor #103: extract validation logic` |
| `perf` | Performance improvement | `perf #67: optimize query with index` |
| `test` | Add or fix tests | `test #88: add unit tests for health` |
| `build` | Build system or dependency changes | `build #120: upgrade gin to v1.10` |
| `ci` | CI/CD changes | `ci #55: add backend test stage` |
| `chore` | Maintenance tasks | `chore #33: update .gitignore` |
| `revert` | Revert a previous commit | `revert #45: revert task scheduler` |

### Issue Reference

The issue reference is **MANDATORY** and uses the `#number` format:

```
feat #123: add email monitoring capability

Closes #123
```

**Rules:**
1. Every commit must reference an existing issue with `#number` in the title
2. If the commit closes the issue, it MUST include `Closes #number` at the end to automatically close the issue
3. A commit can only contain changes for ONE issue
4. If you need to make additional changes, create a new issue first
5. GitLab will automatically close the issue when the commit reaches the main branch

### Description

- Use present imperative: "add" not "added" or "adds"
- First letter in lowercase
- No period at the end
- Maximum 72 characters
- Describe WHAT is done, not HOW

### Complete Examples

```bash
# Simple feature
feat #45: add email monitoring capability

Closes #45

# Fix
fix #78: handle empty request body

Closes #78

# Refactor with explanatory body
refactor #92: split agent entity into smaller components

The Agent entity was becoming too large and handling multiple
responsibilities. This change extracts scheduling logic into
a separate Scheduler value object.

Closes #92

# Breaking change
feat #103!: change response format to JSON:API spec

BREAKING CHANGE: All API responses now follow JSON:API specification.
Clients need to update their parsing logic.

Closes #103

# With co-author (when using AI)
feat #120: implement recurring task system

- Add RecurringTask entity
- Create TaskScheduler service
- Implement cron expression parser

Co-Authored-By: Claude <noreply@anthropic.com>

Closes #120
```

### INCORRECT Examples

```bash
# BAD: No issue reference
feat(api): add new endpoint

# BAD: Mixing changes from multiple issues
feat #45: add monitoring and fix login bug

# BAD: Generic reference without issue
fix: resolve error in handler

# BAD: Old format with parentheses
feat(DEV-45): add new feature
```

### Breaking Changes

For backward-incompatible changes:

1. Add `!` after the type/scope: `feat(api)!: change auth flow`
2. Or include `BREAKING CHANGE:` in the message footer

### Commits with Co-Author

When working with AI assistants, include:

```
Co-Authored-By: Claude <noreply@anthropic.com>
```

---

## Git Workflow

### Branches

| Branch | Purpose |
|--------|---------|
| `main` | Stable production |
| `release` | Pre-production |
| `dev` | Integrated development |
| `predev` | Active development |
| `feature/*` | New features |
| `fix/*` | Bug fixes |
| `hotfix/*` | Urgent production fixes |

### Creating a Feature

```bash
# From predev
git checkout predev
git pull origin predev
git checkout -b feature/DEV-123-short-description

# Work and commit
git add .
git commit -m "feat(scope): description"

# Push and create MR
git push -u origin feature/DEV-123-short-description
```

### Pull Requests / Merge Requests

1. Title must follow Conventional Commits
2. Description must include:
   - Summary of changes
   - How to test
   - Screenshots if applicable
3. Assign reviewers
4. Wait for green CI before merging

---

## Commit Validation

The project has a `commit-msg` hook that validates the format:

```bash
# Install hooks
./scripts/install-git-hooks.sh

# The hook automatically validates each commit
git commit -m "invalid message"  # Error
git commit -m "feat(api): valid message"  # OK
```

### Validation Regex

```regex
^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert) #[0-9]+!?: .{1,100}$
```

The regex validates:
- Valid type (feat, fix, etc.)
- Mandatory issue reference `#XXX`
- Description of up to 100 characters

---

## CI/CD Pipeline

The project uses GitHub Actions for continuous integration. All workflows trigger on push and pull request to `master`.

| Workflow | File | What it does |
|----------|------|-------------|
| **Lint** | `.github/workflows/lint.yml` | Runs `golangci-lint` to enforce code quality |
| **Test** | `.github/workflows/test.yml` | Runs `go test ./... -race` with coverage reporting |
| **Build** | `.github/workflows/build.yml` | Builds binaries for linux/amd64, darwin/amd64, darwin/arm64, windows/amd64 |

### Local Development

Use the Taskfile to run the same checks locally before pushing:

```bash
task lint      # Run golangci-lint
task test      # Run tests with coverage
task build     # Build the binary
task dev       # Start Air for hot-reload
task           # lint + test + build (default)
```

---

## Issue Workflow

- Every issue must use the repository issue templates (Task or Bug Report)
- Every issue that produces code **must include unit tests**
- Reference the issue number in all commits (`feat #N: ...`)
- One issue = one commit = one pull request

---

## Code of Conduct

- Respect the work of other contributors
- Keep discussions technical and constructive
- Follow the project's style guides
- Write tests for new features
- Document significant changes

---

## Checklist before Committing

- [ ] A GitLab issue exists for this change
- [ ] The changes correspond ONLY to that issue
- [ ] The code compiles without errors
- [ ] Tests pass (`task test`)
- [ ] The linter reports no errors (`task lint`)
- [ ] The code is formatted (`task fmt`)
- [ ] The commit message includes the issue reference `#XXX` in the title
- [ ] The commit message includes `Closes #XXX` at the end to close the issue
- [ ] The commit message follows Conventional Commits
- [ ] Tests have been added if it's a new feature
- [ ] Documentation is updated if applicable
