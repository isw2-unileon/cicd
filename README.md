# CI/CD with GitHub Actions

[![CI](https://github.com/isw2-unileon/cicd/actions/workflows/02-continuous-integration.yml/badge.svg)](https://github.com/isw2-unileon/cicd/actions/workflows/02-continuous-integration.yml)
[![Code Quality](https://github.com/isw2-unileon/cicd/actions/workflows/03-code-quality.yml/badge.svg)](https://github.com/isw2-unileon/cicd/actions/workflows/03-code-quality.yml)
[![Coverage](https://github.com/isw2-unileon/cicd/actions/workflows/04-coverage.yml/badge.svg)](https://github.com/isw2-unileon/cicd/actions/workflows/04-coverage.yml)
[![Security](https://github.com/isw2-unileon/cicd/actions/workflows/11-security.yml/badge.svg)](https://github.com/isw2-unileon/cicd/actions/workflows/11-security.yml)
[![Go version](https://img.shields.io/github/go-mod/go-version/isw2-unileon/cicd)](go.mod)
[![License](https://img.shields.io/github/license/isw2-unileon/cicd)](LICENSE)

> Teaching repository for the **Software Engineering II** course — Universidad de León.
> Work through the numbered workflows in order. Each one introduces one new concept.

---

## Part 1 — Why CI/CD?

Imagine a team of 5 developers all pushing code to the same repository. Without any automation:

- Someone pushes code that breaks the build. Nobody knows until another developer tries to pull and compile. Hours of debugging lost.
- A developer forgets to run tests. A bug reaches production. Clients are affected.
- "It works on my machine" — but not on the server, because the server runs a different OS.
- Releasing a new version requires a developer to manually build, copy files, and restart services. Error-prone. Stressful.

**CI/CD solves all of these problems by automating the repetitive, error-prone parts of software delivery.**

### The three terms

| Term | Full name | The idea |
| ---- | --------- | -------- |
| **CI** | Continuous Integration | Every code change is automatically built and tested. Problems are caught in minutes, not days. |
| **CD** | Continuous Delivery | Every passing build is automatically packaged and delivered to a staging environment, ready to deploy. |
| **CD** | Continuous Deployment | Every passing build is automatically deployed to production — no human action required. |

> Most companies practice Continuous Delivery (manual production deploy) rather than full Continuous Deployment. Both are valid.

### The feedback loop

The core metric of CI/CD is **how fast you know if your change broke something**.

```text
Without CI/CD:
  Push code → Forget about it → Review in 2 days → "Oh, this broke staging" → Fix → Repeat

With CI/CD:
  Push code → 3 minutes later: ✅ all green  (or ❌ here is exactly what broke)
```

The shorter the feedback loop, the faster the team can move and the more confident developers are to make changes.

---

## Part 2 — What is GitHub Actions?

GitHub Actions is GitHub's built-in CI/CD platform. When you push code, GitHub reads YAML files from `.github/workflows/` and runs them automatically on virtual machines it provides for free.

**Key vocabulary:**

```text
Workflow    ─── A YAML file in .github/workflows/. Defines what to do and when.
  │
  ├── Trigger (on:)  ─── WHEN the workflow runs: push, pull_request, schedule, manual...
  │
  └── Job            ─── A group of steps that runs on ONE virtual machine.
        │                 Multiple jobs run IN PARALLEL by default.
        │                 Use `needs:` to make them run in sequence.
        │
        └── Step     ─── A single unit of work inside a job.
                         Either a shell command (`run:`) or a pre-built action (`uses:`).
```

**The runner** is the virtual machine GitHub provides. It starts fresh for every run — nothing from a previous run persists. You can choose:

- `ubuntu-latest` — Linux (most common, fastest)
- `macos-latest` — macOS (required for iOS/macOS builds)
- `windows-latest` — Windows (required for Windows-specific testing)

**Actions** are reusable building blocks published by GitHub, companies, or the community. Instead of writing a 20-line script to install Go, you write one line: `uses: actions/setup-go@v5`.

---

## Part 3 — This Repository

This repository contains a deliberately simple Go HTTP API. The application is not the point — it exists to give the workflows something real to build, test, and deploy.

### The application

A REST calculator API with two endpoints:

| Method | Path | Description |
| ------ | ---- | ----------- |
| `GET` | `/health` | Returns `{"status":"ok"}`. Used by load balancers to check the app is alive. |
| `POST` | `/calculate` | Performs a calculation. Body: `{"operation":"add","a":5,"b":3}` |

```bash
# Run locally
make run

# Test it
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"add","a":5,"b":3}'
# → {"result":8}

curl -X POST http://localhost:8080/calculate \
  -d '{"operation":"divide","a":10,"b":0}'
# → {"result":0,"error":"division by zero"}
```

Supported operations: `add`, `subtract`, `multiply`, `divide`.

### Repository structure

```text
.
├── .github/
│   ├── actions/
│   │   └── setup-go-project/      # Custom composite action (Part 4 extra)
│   ├── workflows/
│   │   ├── 01-hello-cicd.yml           ← Start here
│   │   ├── 02-continuous-integration.yml
│   │   ├── 03-code-quality.yml
│   │   ├── 04-coverage.yml
│   │   ├── 05-matrix-builds.yml
│   │   ├── 06-release.yml
│   │   ├── 07-docker.yml
│   │   ├── 08-environments-and-deploy.yml
│   │   ├── 09-manual-workflow.yml
│   │   ├── 10-scheduled.yml
│   │   ├── 11-security.yml
│   │   ├── 12-reusable-workflow.yml
│   │   ├── 13-call-reusable.yml
│   │   └── 14-advanced-features.yml
│   └── dependabot.yml             # Automated dependency updates
├── cmd/server/main.go             # HTTP server
├── internal/calculator/           # Business logic + unit tests
├── Dockerfile                     # Multi-stage container build
├── Makefile                       # Common dev commands
└── go.mod
```

---

## Part 4 — The Workflows

Open each file as you read its section. Every file is **heavily commented** — the comments are part of the lesson.

---

### 01 · Anatomy of a Workflow

**File:** [.github/workflows/01-hello-cicd.yml](.github/workflows/01-hello-cicd.yml)

Before writing any real CI, we need to understand the structure of a workflow file. This workflow does nothing useful — it just prints messages. That is the point.

**What to look at:**

- `name:` — the label shown in the GitHub Actions tab
- `on:` — the trigger. This one uses `push` and `workflow_dispatch` (manual)
- `jobs:` — one job called `hello`
- `runs-on:` — the runner machine (`ubuntu-latest`)
- `steps:` — the list of things to do, in order
- `run:` — executes a shell command
- The built-in environment variables GitHub injects: `$GITHUB_REPOSITORY`, `$GITHUB_SHA`, `$GITHUB_ACTOR`, etc.

**Try it:** Push any commit and watch this workflow run in the Actions tab.

---

### 02 · Continuous Integration

**File:** [.github/workflows/02-continuous-integration.yml](.github/workflows/02-continuous-integration.yml)

This is the most important workflow in the repository. **Every professional project has something like this.**

It runs on every `push` and every `pull_request`. Its job is to answer one question: *does this code compile and do all tests pass?*

**What to look at:**

- `actions/checkout@v4` — clones the repository onto the runner. Without this, the runner has no code.
- `actions/setup-go@v5` — installs Go. The `go-version-file: go.mod` option reads the version from `go.mod` automatically.
- `go mod verify` — checks that nobody tampered with the dependencies since `go.sum` was written.
- `go build ./...` — compiles. If there is a syntax error, the workflow fails here.
- `go test -race ./...` — runs all tests. The `-race` flag enables the data race detector.

**The key concept — branch protection rules:**

Go to *Settings → Branches → Add rule* and require this workflow to pass before a PR can be merged. This means **it is physically impossible to merge broken code**. The CI is the gatekeeper.

```text
Developer opens PR
        ↓
  This workflow runs automatically
        ↓
  ❌ Tests fail → PR blocked, cannot merge
  ✅ Tests pass → PR can be reviewed and merged
```

---

### 03 · Code Quality

**File:** [.github/workflows/03-code-quality.yml](.github/workflows/03-code-quality.yml)

CI verifies correctness. This workflow verifies **quality**. It runs three independent jobs in parallel:

| Job | Tool | What it catches |
| --- | ---- | --------------- |
| `vet` | `go vet` | Suspicious constructs: wrong `Printf` format, unreachable code, etc. |
| `format` | `gofmt` | Files that are not properly formatted. The workflow fails if any exist. |
| `lint` | `golangci-lint` | 10+ linters at once: unused variables, unchecked errors, style issues, and more. |

**What to look at:**

- The `paths:` filter on the trigger — this workflow only runs when `.go` files change. There is no point running the linter if only the README was updated.
- The three jobs run **in parallel** — GitHub starts all three at the same time. This is faster than running them sequentially.
- The `golangci-lint-action` — it caches the linter binary and results, making repeated runs much faster.

**The key concept — fail fast:**

Short, cheap checks (vet, fmt) run before expensive ones (lint). If the code is not even formatted, there is no need to run the full linter suite.

---

### 04 · Test Coverage

**File:** [.github/workflows/04-coverage.yml](.github/workflows/04-coverage.yml)

This workflow measures what percentage of the code is executed by the test suite. Low coverage does not mean the code is bad — but it tells you where tests are missing.

**What to look at:**

- `go test -coverprofile=coverage.out` — runs tests and writes coverage data to a file
- `go tool cover -func=coverage.out` — prints a per-function summary in the terminal
- `go tool cover -html=coverage.out` — generates an HTML report where you can see exactly which lines are covered (green) and which are not (red)
- `actions/upload-artifact@v4` — saves files produced by the job so team members can download them after the run. The HTML report is stored here.
- `$GITHUB_STEP_SUMMARY` — a special file. Whatever you write to it appears as formatted markdown in the workflow summary page. Students should see the coverage table appear there after the run.

---

### 05 · Matrix Builds

**File:** [.github/workflows/05-matrix-builds.yml](.github/workflows/05-matrix-builds.yml)

**Problem:** Your code works on your laptop (macOS, Go 1.24). Will it work on the production server (Linux, Go 1.26)? On a colleague's Windows machine?

**Solution:** Test on all combinations automatically.

```yaml
matrix:
  go-version: ['1.24', '1.25', '1.26']
  os: [ubuntu-latest, macos-latest, windows-latest]
```

This generates **9 jobs** (3 versions × 3 OS) that all run in parallel. If your code has a platform-specific bug, one of them will catch it.

**What to look at:**

- `strategy.matrix` — the definition of the combinations
- `strategy.fail-fast: false` — do not cancel the remaining 8 jobs if one fails. Let them all finish so you can see the full picture.
- `${{ matrix.go-version }}` and `${{ matrix.os }}` — how to reference matrix variables in the job
- `include:` — add extra variables to specific combinations. Used here to run coverage only on the "primary" combination (Ubuntu + latest Go) to avoid uploading 9 identical reports.

---

### 06 · Automated Releases

**File:** [.github/workflows/06-release.yml](.github/workflows/06-release.yml)

Before CI/CD, releasing software meant: build locally, zip files, upload to a server, update a wiki. This workflow does all of it automatically when a developer pushes a version tag.

```bash
# A developer does this locally:
git tag v1.2.3
git push --tags

# GitHub Actions does the rest automatically:
# → Compiles for Linux, macOS, Windows (both Intel and ARM)
# → Creates a GitHub Release page
# → Uploads all 5 binaries as downloadable assets
```

**What to look at:**

- `on: push: tags: ['v*.*.*']` — this workflow only triggers on version tags, never on regular commits
- `permissions: contents: write` — workflows need explicit permission to create releases. By default they are read-only.
- `GOOS` and `GOARCH` environment variables — Go's cross-compilation. From a single Linux runner, you can produce native binaries for every platform. No need for a Mac to build a Mac binary.
- `gh release create` — GitHub CLI is pre-installed on all runners. Used to create the release and upload the assets.

---

### 07 · Docker

**File:** [.github/workflows/07-docker.yml](.github/workflows/07-docker.yml)

Containers solve the "works on my machine" problem permanently — the entire runtime environment ships with the application. This workflow builds a Docker image and pushes it to GitHub Container Registry (`ghcr.io`), which is free for public repositories.

**What to look at:**

- `docker/metadata-action` — automatically computes image tags based on the trigger. Push to `main` → tags `latest` and `main`. Push tag `v1.2.3` → tags `1.2.3`, `1.2`, `1`. No manual tag management.
- `docker/setup-qemu-action` — QEMU emulates other CPU architectures. This allows building an ARM64 image from an AMD64 runner (important for Apple Silicon and AWS Graviton).
- `cache-from: type=gha` — Docker layer caching. Unchanged layers are reused from previous runs. Saves minutes on every build.
- `push: ${{ github.event_name != 'pull_request' }}` — pull requests build the image to verify it compiles, but do not push it to the registry. Forked PRs do not have access to secrets.
- Also look at the [Dockerfile](Dockerfile) — it uses a **multi-stage build**: stage 1 compiles the binary with the full Go toolchain; stage 2 copies only the binary into a minimal `scratch` image. Result: ~10 MB instead of ~800 MB.

---

### 08 · Environments and Deployment

**File:** [.github/workflows/08-environments-and-deploy.yml](.github/workflows/08-environments-and-deploy.yml)

A real deployment pipeline does not go directly from code to production. It goes through stages:

```text
Tests pass → Build artifact → Deploy to staging → Human approves → Deploy to production
```

**GitHub Environments** are named deployment targets (`staging`, `production`) with configurable rules:

- **Required reviewers:** The workflow pauses and sends a notification. A designated person must click "Approve" before the production job runs.
- **Wait timer:** Add a delay before deployment (e.g. 10 minutes to allow monitoring alerts to fire).
- **Allowed branches:** Only `main` can deploy to production.

**What to look at:**

- The `needs:` chain: `test` → `build` → `deploy-staging` → `deploy-production`. Jobs run in order, each waiting for the previous to succeed.
- `environment: production` — links the job to the GitHub Environment. GitHub shows a deployment badge on the repository home page and tracks the history.
- `actions/upload-artifact` and `actions/download-artifact` — the compiled binary is built once in the `build` job and downloaded by both deploy jobs. No recompilation.

**Setup required before this works:**

Go to *Settings → Environments*, create `staging` and `production`, then add yourself as a required reviewer for `production`.

---

### 09 · Manual Workflows

**File:** [.github/workflows/09-manual-workflow.yml](.github/workflows/09-manual-workflow.yml)

Not every workflow should run automatically. Some tasks — database migrations, cache invalidation, deploying a specific version — should be triggered by a human, with explicit parameters.

`workflow_dispatch` adds a "Run workflow" button in the GitHub Actions tab. You fill in a form and click the button.

**What to look at:**

- The four input types: `choice` (dropdown), `string` (text), `boolean` (checkbox), `number`
- `${{ inputs.environment }}` — how to read the input value inside the workflow
- `if: inputs.run_migrations == true` — conditional step that only runs when the checkbox was checked
- `environment: ${{ inputs.environment }}` — dynamically targets the environment the user selected

---

### 10 · Scheduled Jobs

**File:** [.github/workflows/10-scheduled.yml](.github/workflows/10-scheduled.yml)

Some tasks should run on a clock, not triggered by code changes: nightly full test runs to catch flaky tests, weekly checks for new security vulnerabilities, monthly dependency reports.

GitHub Actions supports standard Unix **cron syntax**:

```text
┌─ minute   (0-59)
│ ┌─ hour   (0-23, UTC)
│ │ ┌─ day  (1-31)
│ │ │ ┌─ month   (1-12)
│ │ │ │ ┌─ weekday (0-7, 0=Sunday)
│ │ │ │ │
0 2 * * *     → Every day at 02:00 UTC
0 9 * * 1     → Every Monday at 09:00 UTC
*/15 * * * *  → Every 15 minutes
```

> **Important:** All GitHub Actions cron runs in **UTC**. If your team is in UTC+2, a "9am" job should be scheduled at `0 7 * * *`.

**What to look at:**

- Two separate schedules in the same workflow (nightly tests + weekly dependency check)
- `if: github.event.schedule == '0 2 * * *'` — how to run different jobs depending on which schedule triggered the workflow
- `govulncheck` — checks for known security vulnerabilities in Go dependencies

---

### 11 · Security Scanning

**File:** [.github/workflows/11-security.yml](.github/workflows/11-security.yml)

Security should be part of the development process, not something you think about after a breach. The term **"shift left"** means moving security checks earlier in the pipeline — catching vulnerabilities when the code is written, not after it is deployed.

This workflow has two independent security jobs:

**CodeQL** — GitHub's static analysis engine. It reads your source code, builds a model of how data flows through the program, and detects vulnerabilities (SQL injection, path traversal, use of dangerous functions, etc.). Results appear in the *Security* tab of the repository.

**govulncheck** — checks whether your Go dependencies have known CVEs (Common Vulnerabilities and Exposures). Crucially, it only alerts you if your code **actually calls** the vulnerable function — not just if you depend on a vulnerable version. This dramatically reduces false positives.

**What to look at:**

- `permissions: security-events: write` — required to upload CodeQL results to the Security tab
- `github/codeql-action/init` → `autobuild` → `analyze` — the three-step CodeQL process
- `repo-checkout: false` on `govulncheck-action` — the action runs its own internal checkout by default. Since we already checked out in the previous step, we disable the duplicate to avoid an authentication conflict.

---

### 12 & 13 · Reusable Workflows

**Files:** [.github/workflows/12-reusable-workflow.yml](.github/workflows/12-reusable-workflow.yml), [.github/workflows/13-call-reusable.yml](.github/workflows/13-call-reusable.yml)

Imagine you have 5 repositories and each one has a CI workflow. They all do the same thing: checkout, setup Go, test, lint. When you need to update them, you update 5 files. When one drifts from the others, you get inconsistency.

**Reusable workflows** solve this. Define the logic once, call it from everywhere.

`workflow_call` is the trigger that makes a workflow callable:

```yaml
# 12-reusable-workflow.yml — defines the logic
on:
  workflow_call:
    inputs:
      go-version:
        type: string
    outputs:
      coverage:
        value: ${{ jobs.test.outputs.coverage }}

# 13-call-reusable.yml — uses the logic
jobs:
  ci:
    uses: ./.github/workflows/12-reusable-workflow.yml
    with:
      go-version: '1.24'
    secrets: inherit
```

**What to look at:**

- `inputs:` — typed parameters the caller passes in (like function arguments)
- `outputs:` — values the caller can read after the workflow finishes (like return values)
- `secrets: inherit` — passes all secrets from the calling workflow automatically, without listing them by name
- In `13-call-reusable.yml`: the `report` job uses `needs.ci.outputs.coverage` to read the output produced by the reusable workflow

---

### 14 · Advanced Features

**File:** [.github/workflows/14-advanced-features.yml](.github/workflows/14-advanced-features.yml)

A collection of features that make workflows more robust and efficient in production.

**Concurrency control** — without this, pushing 3 commits quickly starts 3 workflow runs simultaneously. They all test the same branch but the first two results are irrelevant by the time they finish. `concurrency` ensures only the latest run is active:

```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true   # cancel the old run when a new one starts
```

**Timeouts** — a job that gets stuck (waiting for a service, infinite loop in a test) will run until GitHub's 6-hour limit. `timeout-minutes: 10` kills it after 10 minutes, freeing the runner and alerting you quickly.

**Job outputs** — the mechanism for passing data between jobs. A step writes to `$GITHUB_OUTPUT`; the job exposes it as an `output:`; the next job reads it via `needs.jobname.outputs.key`.

**Conditional steps** — `if:` accepts any GitHub Actions expression. Common patterns:

```yaml
if: github.ref == 'refs/heads/main'          # only on main branch
if: github.event_name == 'pull_request'       # only on PRs
if: always()                                  # run even if previous step failed
if: failure()                                 # only run if something above failed
```

---

## Part 5 — Custom Actions

**File:** [.github/actions/setup-go-project/action.yml](.github/actions/setup-go-project/action.yml)

Actions come in three flavours:

| Type | How it works | Best for |
| ---- | ------------ | -------- |
| **JavaScript** | Runs a Node.js script | Complex logic, cross-platform |
| **Docker** | Runs a container | Specific environment requirements |
| **Composite** | Chains multiple steps | Reusing a group of steps |

The action in this repository is **composite** — it bundles the 3 setup steps every workflow repeats (checkout, setup-go, verify deps) into a single call:

```yaml
# Instead of 3 separate steps in every workflow, you write one:
- name: Setup project
  uses: ./.github/actions/setup-go-project
  with:
    go-version: '1.24'
```

**The difference between a custom action and a reusable workflow:**

- A **reusable workflow** replaces an entire job (it runs on its own runner).
- A **custom action** is a step inside a job (it runs on the caller's runner).

---

## Part 6 — Dependabot

**File:** [.github/dependabot.yml](.github/dependabot.yml)

Dependencies have vulnerabilities. New versions ship with security fixes. Keeping dependencies up to date manually is tedious and easy to forget.

**Dependabot** is a GitHub bot that automatically opens pull requests when newer versions of your dependencies are available. Your CI runs on those PRs just like any other — if the update breaks something, the tests fail and you do not merge.

This configuration watches two ecosystems:

- **Go modules** (`go.mod`) — checks for updated packages
- **GitHub Actions** — checks for updated action versions (`actions/checkout@v4` → `v5`, etc.)

> Most teams configure Dependabot to group minor and patch updates into a single weekly PR, and create separate PRs for major updates that might have breaking changes.

---

## Part 7 — Reference

### Local development commands

```bash
make run        # Start server on :8080
make test       # Run all tests
make coverage   # Tests + HTML coverage report
make lint       # Run golangci-lint
make fmt        # Format all .go files
make vet        # Run go vet
make build      # Compile to bin/server
make clean      # Remove build artifacts
```

### GitHub Contexts

Inside any workflow file, `${{ }}` expressions give you access to information about the current run:

| Expression | Example value | Meaning |
| ---------- | ------------- | ------- |
| `github.event_name` | `push` | What triggered this run |
| `github.ref` | `refs/heads/main` | Full git ref |
| `github.ref_name` | `main` | Short branch or tag name |
| `github.sha` | `a1b2c3d4...` | Commit SHA |
| `github.actor` | `jferrl` | Who triggered the run |
| `github.repository` | `isw2-unileon/cicd` | Owner/repo |
| `github.run_number` | `42` | Increments with every run |
| `secrets.NAME` | *(hidden)* | An encrypted repository secret |
| `env.NAME` | `calculator-api` | An environment variable |
| `inputs.NAME` | `staging` | A `workflow_dispatch` input |
| `needs.JOB.outputs.KEY` | `v1.2.3` | Output from a previous job |

### All workflow triggers

| Trigger | When it fires |
| ------- | ------------- |
| `push` | On every commit pushed |
| `pull_request` | When a PR is opened, updated, or reopened |
| `workflow_dispatch` | Manual trigger from the Actions tab |
| `schedule` | On a cron schedule |
| `workflow_call` | Called by another workflow |
| `release` | When a GitHub Release is created/published |
| `issue_comment` | When someone comments on an issue or PR |
| `workflow_run` | When another workflow finishes |
