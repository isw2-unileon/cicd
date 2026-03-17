# CI/CD with GitHub Actions

[![CI](https://github.com/isw2-unileon/cicd/actions/workflows/02-continuous-integration.yml/badge.svg)](https://github.com/isw2-unileon/cicd/actions/workflows/02-continuous-integration.yml)
[![Code Quality](https://github.com/isw2-unileon/cicd/actions/workflows/03-code-quality.yml/badge.svg)](https://github.com/isw2-unileon/cicd/actions/workflows/03-code-quality.yml)
[![Coverage](https://github.com/isw2-unileon/cicd/actions/workflows/04-coverage.yml/badge.svg)](https://github.com/isw2-unileon/cicd/actions/workflows/04-coverage.yml)
[![Security](https://github.com/isw2-unileon/cicd/actions/workflows/11-security.yml/badge.svg)](https://github.com/isw2-unileon/cicd/actions/workflows/11-security.yml)
[![Go version](https://img.shields.io/github/go-mod/go-version/isw2-unileon/cicd)](go.mod)
[![License](https://img.shields.io/github/license/isw2-unileon/cicd)](LICENSE)

A teaching repository for the **Software Engineering II** course. This repo demonstrates CI/CD principles and GitHub Actions from the ground up, using a simple Go HTTP API as the application under delivery.

---

## What is CI/CD?

| Term | Full name | What it means |
| ---- | --------- | ------------- |
| **CI** | Continuous Integration | Automatically build and test every code change |
| **CD** | Continuous Delivery | Automatically deliver tested code to a staging environment |
| **CD** | Continuous Deployment | Automatically deploy to production without human intervention |

The goal: **shorten the feedback loop** between writing code and knowing it works.

```text
Developer pushes code
        ↓
  GitHub receives push
        ↓
  Workflow triggers (automatic)
        ↓
  ┌─────────────────────────────┐
  │  CI: Build → Test → Lint   │  ← Seconds to minutes
  └─────────────────────────────┘
        ↓ (if all green)
  ┌─────────────────────────────┐
  │  CD: Build image → Deploy  │  ← Minutes
  └─────────────────────────────┘
        ↓
  Code is live in production
```

---

## Repository Structure

```text
.
├── .github/
│   ├── actions/
│   │   └── setup-go-project/     # Custom composite action
│   │       └── action.yml
│   ├── workflows/
│   │   ├── 01-hello-cicd.yml           # Anatomy of a workflow
│   │   ├── 02-continuous-integration.yml # Build + test on every push
│   │   ├── 03-code-quality.yml         # Linting and formatting
│   │   ├── 04-coverage.yml             # Test coverage reports
│   │   ├── 05-matrix-builds.yml        # Test on multiple OS/Go versions
│   │   ├── 06-release.yml              # Automated GitHub releases
│   │   ├── 07-docker.yml               # Docker build & push to ghcr.io
│   │   ├── 08-environments-and-deploy.yml # Staging/production environments
│   │   ├── 09-manual-workflow.yml      # Manual triggers with inputs
│   │   ├── 10-scheduled.yml            # Cron/scheduled jobs
│   │   ├── 11-security.yml             # CodeQL + vulnerability scanning
│   │   ├── 12-reusable-workflow.yml    # Reusable workflow (called by 13)
│   │   ├── 13-call-reusable.yml        # Calling a reusable workflow
│   │   └── 14-advanced-features.yml    # Caching, concurrency, outputs
│   └── dependabot.yml                  # Automated dependency updates
├── cmd/server/main.go                  # HTTP server entry point
├── internal/calculator/
│   ├── calculator.go                   # Business logic
│   └── calculator_test.go              # Unit tests
├── Dockerfile                          # Multi-stage Docker build
├── Makefile                            # Common development commands
├── go.mod
└── .golangci.yml                       # Linter configuration
```

---

## Workflows — Learning Path

Work through the workflows in order. Each one introduces a new concept.

### 01 · Hello, CI/CD

**File:** [.github/workflows/01-hello-cicd.yml](.github/workflows/01-hello-cicd.yml)

The simplest possible workflow. Teaches the basic structure:

- `name:` — displayed in the GitHub Actions tab
- `on:` — triggers (push, pull_request, workflow_dispatch, schedule, ...)
- `jobs:` — groups of steps running on the same machine
- `steps:` — individual commands (`run:`) or pre-built actions (`uses:`)
- `runs-on:` — the virtual machine GitHub provides

### 02 · Continuous Integration

**File:** [.github/workflows/02-continuous-integration.yml](.github/workflows/02-continuous-integration.yml)

The foundation of every CI pipeline. On every push and pull request:

1. Clone the repo (`actions/checkout`)
2. Install Go (`actions/setup-go`)
3. Verify dependencies (`go mod verify`)
4. Compile (`go build ./...`)
5. Run tests (`go test -race ./...`)

> **Key insight:** Pull request checks — configure branch protection rules so a PR *cannot be merged* until this workflow passes.

### 03 · Code Quality

**File:** [.github/workflows/03-code-quality.yml](.github/workflows/03-code-quality.yml)

Enforce code standards automatically:

- `go vet` — built-in static analysis
- `gofmt` — fail if code is not properly formatted
- `golangci-lint` — 50+ linters in one tool

> **Key insight:** Path filters — only run when `.go` files change. No point linting if only the README was updated.

### 04 · Test Coverage

**File:** [.github/workflows/04-coverage.yml](.github/workflows/04-coverage.yml)

- Generate coverage profiles with `go test -coverprofile`
- Upload HTML reports as **workflow artifacts** (downloadable from the UI)
- Write custom markdown to the **job summary** using `$GITHUB_STEP_SUMMARY`

### 05 · Matrix Builds

**File:** [.github/workflows/05-matrix-builds.yml](.github/workflows/05-matrix-builds.yml)

Test across multiple Go versions (1.21, 1.22, 1.23) and operating systems (Linux, macOS, Windows) **in parallel** using a single job definition.

```yaml
strategy:
  matrix:
    go-version: ['1.21', '1.22', '1.23']
    os: [ubuntu-latest, macos-latest, windows-latest]
```

> **Key insight:** This catches "works on my machine" problems before they reach production.

### 06 · Automated Releases

**File:** [.github/workflows/06-release.yml](.github/workflows/06-release.yml)

Triggered when a developer pushes a version tag (`git tag v1.2.3 && git push --tags`):

1. Cross-compile for Linux, macOS, Windows (amd64 and arm64)
2. Create a GitHub Release
3. Upload all binaries as downloadable release assets

> **Key insight:** Cross-compilation — Go can produce native binaries for every platform from a single Linux runner using `GOOS` and `GOARCH` environment variables.

### 07 · Docker

**File:** [.github/workflows/07-docker.yml](.github/workflows/07-docker.yml)

- Build a multi-platform Docker image (linux/amd64 + linux/arm64)
- Push to GitHub Container Registry (`ghcr.io`)
- Smart image tagging via `docker/metadata-action`
- Layer caching with `cache-from: type=gha`

### 08 · Environments & Deployment

**File:** [.github/workflows/08-environments-and-deploy.yml](.github/workflows/08-environments-and-deploy.yml)

A complete deployment pipeline:

```text
test → build → deploy-staging → (manual approval) → deploy-production
```

- GitHub **Environments** (`staging`, `production`) with protection rules
- The production job **pauses** and sends a notification asking for approval
- Environment-specific secrets
- Deployment status shown on the repository home page

> **Setup required:** Create environments in *Settings → Environments* and add required reviewers to `production`.

### 09 · Manual Workflows

**File:** [.github/workflows/09-manual-workflow.yml](.github/workflows/09-manual-workflow.yml)

`workflow_dispatch` lets you trigger a workflow manually from the GitHub UI with typed inputs:

- `choice` — dropdown
- `string` — text field
- `boolean` — checkbox
- `number` — numeric field

### 10 · Scheduled Jobs

**File:** [.github/workflows/10-scheduled.yml](.github/workflows/10-scheduled.yml)

Run workflows on a cron schedule. Standard Unix cron syntax:

```text
┌─ minute (0-59)
│ ┌─ hour (0-23)
│ │ ┌─ day of month (1-31)
│ │ │ ┌─ month (1-12)
│ │ │ │ ┌─ day of week (0-7, 0=Sunday)
│ │ │ │ │
* * * * *
```

Examples: nightly test runs, weekly dependency checks, monthly reports.

> **Note:** All GitHub Actions cron runs in **UTC**.

### 11 · Security Scanning

**File:** [.github/workflows/11-security.yml](.github/workflows/11-security.yml)

"Shift left" on security — find vulnerabilities during development, not after deployment:

- **CodeQL** — GitHub's static analysis engine, results appear in the Security tab
- **govulncheck** — scans Go dependencies for known CVEs, only alerts when you actually *call* the vulnerable function

### 12 & 13 · Reusable Workflows

**Files:** [.github/workflows/12-reusable-workflow.yml](.github/workflows/12-reusable-workflow.yml), [.github/workflows/13-call-reusable.yml](.github/workflows/13-call-reusable.yml)

Define a workflow once (`workflow_call`) and call it from multiple other workflows. Like a function call but for entire jobs:

```yaml
# In the calling workflow:
jobs:
  ci:
    uses: ./.github/workflows/12-reusable-workflow.yml
    with:
      go-version: '1.23'
    secrets: inherit
```

### 14 · Advanced Features

**File:** [.github/workflows/14-advanced-features.yml](.github/workflows/14-advanced-features.yml)

| Feature | What it does |
| ------- | ------------ |
| `concurrency:` | Cancel in-progress runs when a new commit is pushed |
| `timeout-minutes:` | Kill jobs that run too long |
| `if:` expressions | Conditional steps based on context (branch, event type, etc.) |
| Job outputs | Pass data between jobs via `$GITHUB_OUTPUT` |
| `$GITHUB_STEP_SUMMARY` | Write rich markdown to the workflow summary page |

---

## Custom Action

**File:** [.github/actions/setup-go-project/action.yml](.github/actions/setup-go-project/action.yml)

A **composite action** that bundles the 3 setup steps every workflow repeats (checkout, setup-go, verify deps) into a single reusable call:

```yaml
- name: Setup project
  uses: ./.github/actions/setup-go-project
  with:
    go-version: '1.23'
```

---

## The Application

A minimal REST API that exposes basic arithmetic operations. It exists only to give the CI/CD workflows something real to build and test.

### Endpoints

| Method | Path | Description |
| ------ | ---- | ----------- |
| GET | `/health` | Health check — returns `{"status":"ok"}` |
| POST | `/calculate` | Perform a calculation |

### Calculate Example

```bash
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"add","a":5,"b":3}'
# → {"result":8}

curl -X POST http://localhost:8080/calculate \
  -d '{"operation":"divide","a":10,"b":0}'
# → {"result":0,"error":"division by zero"}
```

Supported operations: `add`, `subtract`, `multiply`, `divide`.

---

## Local Development

```bash
# Run the server
make run

# Run all tests
make test

# Run tests with coverage report
make coverage

# Format code
make fmt

# Run linter
make lint

# Build binary
make build
```

---

## GitHub Contexts Reference

Inside workflow files, you can access information about the run using **contexts**:

| Context | Example | Description |
| ------- | ------- | ----------- |
| `github.event_name` | `push` | What triggered the workflow |
| `github.ref` | `refs/heads/main` | Full git ref |
| `github.ref_name` | `main` | Short branch/tag name |
| `github.sha` | `a1b2c3d4` | Commit SHA |
| `github.actor` | `username` | Who triggered the run |
| `github.repository` | `org/repo` | Repository name |
| `github.run_number` | `42` | Auto-incrementing run counter |
| `secrets.MY_SECRET` | — | Repository/environment secret |
| `env.MY_VAR` | — | Environment variable |
| `inputs.my_input` | — | `workflow_dispatch` input |
| `needs.job_id.outputs.x` | — | Output from a previous job |

---

## Key GitHub Actions Concepts Summary

```text
Workflow              A YAML file in .github/workflows/
  └── Trigger (on:)  push / pull_request / schedule / workflow_dispatch / ...
  └── Job            Runs on a fresh virtual machine (runner)
        └── Step     Either `run:` (shell command) or `uses:` (action)

Runner               GitHub-hosted VM: ubuntu-latest, macos-latest, windows-latest
Action               Reusable step: from Marketplace, another repo, or local .github/actions/
Artifact             File(s) produced by a job, downloadable from the UI
Environment          Named deployment target (staging, production) with protection rules
Secret               Encrypted value, never shown in logs
Context              Read-only data about the run (github.*, env.*, secrets.*, ...)
Expression           ${{ }} syntax for dynamic values in workflow files
```
