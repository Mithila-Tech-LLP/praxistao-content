# Chapter 92: CI/CD for a Blockchain Project

GoChain now spans thirteen volumes of packages — `crypto`, `core`, `consensus`, `wallet`, `network`, `storage`, `vm`, `api`, `metrics` — each with its own test suite, plus a `Dockerfile` (Chapter 86) that packages all of it into a deployable image. Running every test and rebuilding that image by hand before every deploy does not scale, and worse, it is exactly the kind of tedious manual step a tired human eventually skips. This chapter wires up GitHub Actions so every push runs the full test suite automatically, and every tagged release builds and publishes a new Docker image with zero manual steps.

## Table of Contents

1. [Why Automate Deploys for a Blockchain Project Specifically](#1-why-automate-deploys-for-a-blockchain-project-specifically)
2. [Continuous Integration vs. Continuous Deployment](#2-continuous-integration-vs-continuous-deployment)
3. [GitHub Actions Concepts](#3-github-actions-concepts)
4. [What "The Full Test Suite" Means for GoChain](#4-what-the-full-test-suite-means-for-gochain)
5. [Writing ci.yml — the Test Job](#5-writing-ciyml--the-test-job)
6. [Writing ci.yml — the Build Job](#6-writing-ciyml--the-build-job)
7. [Container Registries and GHCR](#7-container-registries-and-ghcr)
8. [Writing release.yml — Tag, Build, Push](#8-writing-releaseyml--tag-build-push)
9. [From Registry to Running Node — the Deploy Step](#9-from-registry-to-running-node--the-deploy-step)
10. [Secrets: Keeping Credentials Out of the Repo](#10-secrets-keeping-credentials-out-of-the-repo)
11. [Watching a Pipeline Run](#11-watching-a-pipeline-run)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. Why Automate Deploys for a Blockchain Project Specifically

In most web projects, a bad deploy means a page returns a 500 error until someone rolls it back. In a blockchain project, the stakes of a bad deploy are subtly higher: a node running buggy consensus code (Chapter 24-27) can validate blocks it should have rejected, or a broken transaction-signing path (Chapter 33) can silently corrupt the mempool other honest nodes are gossiping from (Chapter 48). Because every node is independently trusted to enforce the exact same rules, "deploy first, notice the bug later" is far more expensive here than in most software — by the time you notice, the buggy node may have already propagated bad data to peers.

The single best defense against shipping a broken node is refusing to let *any* change reach production without first passing the same test suite you already spent twelve volumes writing. **Continuous Integration** (CI) makes that non-negotiable and automatic: nobody can forget to run `go test ./...` before merging, because a robot runs it for them, every single time, on every single push.

---

## 2. Continuous Integration vs. Continuous Deployment

**Continuous Integration (CI)** is the practice of automatically building and testing every change the moment it is pushed, so problems are caught in minutes rather than discovered after they have already been merged and deployed. **Continuous Deployment (CD)** goes one step further: once a change passes CI and meets some release criterion (in GoChain's case, being tagged), it is automatically shipped to production with no manual "click deploy" step at all.

```
   CI (every push)                       CD (on a tag)
   ----------------                      --------------
   push code                              push a tag  (v1.4.0)
       |                                       |
       v                                       v
   run full test suite                    run full test suite
       |                                       |
       v                                       v
   build Docker image                      build Docker image
   (verify it builds)                          |
                                                v
                                          push image to registry
                                                |
                                                v
                                          deploy to the cloud VM
```

Both halves matter for GoChain: CI on every push keeps `main` always in a known-good state; CD on tags turns "ship version 1.4.0" into "type one `git tag` command" instead of a manual SSH-and-rebuild session on the VM from Chapter 88.

---

## 3. GitHub Actions Concepts

**GitHub Actions** is GitHub's built-in automation system: YAML files in `.github/workflows/` describe **workflows** that run in response to repository events. A workflow is made of one or more **jobs**, each of which runs on a fresh, isolated virtual machine called a **runner**. Each job is a sequence of **steps** — either a shell command or a reusable **action** (a pre-built step someone else published, referenced like `actions/checkout@v4`).

```
  .github/workflows/ci.yml
        |
        | triggers on: push, pull_request
        v
  +---------------------------------------------+
  |  job: test                                    |
  |  runs-on: ubuntu-latest (a fresh runner)      |
  |  steps:                                       |
  |    1. actions/checkout@v4                     |
  |    2. actions/setup-go@v5                     |
  |    3. go test ./...                           |
  +---------------------------------------------+
  |  job: build   (needs: test)                   |
  |  steps:                                       |
  |    1. actions/checkout@v4                     |
  |    2. docker build -t gochain:ci .             |
  +---------------------------------------------+
```

Every job runs on a brand-new, disposable machine with nothing pre-installed except a base OS image — this is why every workflow starts with `actions/checkout` (to pull your code onto that fresh machine) and `actions/setup-go` (to install the exact Go toolchain version your project needs). Nothing about your own laptop's environment leaks in, which is exactly the point: if a workflow passes, it proves the project builds and tests cleanly from nothing but the repository itself, the same guarantee anyone cloning it fresh would have.

---

## 4. What "The Full Test Suite" Means for GoChain

By Chapter 91, GoChain's module contains tests spread across every package written since Volume 1: `crypto` (Chapter 09, 13), `core` (Chapters 17-19, 32-35), `consensus` (Chapters 25-27, 77), `wallet` (Chapters 39-40), `network` (Chapters 44-51), `storage` (Chapters 54-58), `vm` (Chapters 62-68), and `api` (Chapters 70-72). Go's tooling makes "run every one of those" a single command, because `go test` recursively discovers every `_test.go` file under a package path:

```bash
go test ./...
```

The `./...` pattern means "this directory and every package nested under it, recursively" — one flag, and CI does not need a growing, hand-maintained list of package paths as the course (and the module) keeps growing. Two extra flags make the CI run more useful than a bare pass/fail:

```bash
go test ./... -race -cover
```

`-race` enables Go's built-in **race detector**, which catches a specific and notoriously hard-to-reproduce class of bug: two goroutines accessing the same memory without proper synchronization. GoChain is full of concurrency by design — the miner (Chapter 27), the P2P node (Chapter 46), the mempool (Chapter 34) — so a race-free result on every push is a meaningful, ongoing guarantee, not a one-time check. `-cover` reports what percentage of each package's code the tests actually exercised, surfacing gaps before they turn into unnoticed production bugs.

---

## 5. Writing ci.yml — the Test Job

```yaml
# .github/workflows/ci.yml
#
# Runs on every push and every pull request, so a broken change is caught
# the moment it exists — never merged first and discovered later.
name: CI

on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      # Pulls the exact commit that triggered this workflow onto the
      # runner's fresh filesystem — without this step, there is no code
      # here at all.
      - name: Checkout code
        uses: actions/checkout@v4

      # Installs the Go toolchain version pinned in go.mod, so CI always
      # builds with the same compiler version contributors use locally.
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version-file: "go.mod"
          cache: true # caches downloaded modules between runs, speeding up CI

      # go vet catches suspicious constructs (wrong Printf verbs, unreachable
      # code) that compile fine but are almost always bugs.
      - name: Vet
        run: go vet ./...

      # The full suite from every volume, with the race detector and
      # coverage reporting from Section 4.
      - name: Test
        run: go test ./... -race -cover -v

      # gofmt -l lists any file that is not correctly formatted; a
      # non-empty result should fail the build, since inconsistent
      # formatting is exactly the kind of thing a robot should enforce
      # instead of a human reviewer nagging about it in every PR.
      - name: Check formatting
        run: |
          UNFORMATTED=$(gofmt -l .)
          if [ -n "$UNFORMATTED" ]; then
            echo "The following files are not gofmt-formatted:"
            echo "$UNFORMATTED"
            exit 1
          fi
```

`go-version-file: "go.mod"` reads the `go 1.22` (or whichever version) line already declared in GoChain's `go.mod` from Chapter 06, so the CI Go version and the project's declared minimum version can never silently drift apart. The formatting check's shell block is a small but real example of "why config, not just code": it turns a style convention (Chapter 07 onward, every file has been `gofmt`-clean) into an enforced rule nobody can accidentally violate on `main`.

---

## 6. Writing ci.yml — the Build Job

The build job proves the `Dockerfile` from Chapter 86 still produces a working image — this catches a real, common failure mode where tests pass but the Docker build itself is broken (a missing file in the build context, a stale dependency in the image's base layer). It only runs after `test` succeeds, using `needs`, so a broken test suite never wastes runner time building an image nobody should trust yet:

```yaml
  build:
    name: Build Docker image
    runs-on: ubuntu-latest
    needs: test # only build if every test passed
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      # Builds the exact multi-stage Dockerfile from Chapter 86. Tagging
      # it "ci" (not a real version) makes clear this image is a build
      # verification artifact, not something meant to be deployed.
      - name: Build image
        run: docker build -t gochain:ci .

      # A minimal smoke test: start the just-built image and confirm the
      # binary at least runs and prints its version, rather than crashing
      # immediately on startup — a cheap check that catches an entire
      # class of "it built, but it's broken" failures.
      - name: Smoke test the image
        run: |
          docker run --rm gochain:ci gochain --version
```

This job runs on every push, not just tagged releases — the whole point is catching a broken `Dockerfile` or a broken binary *before* anyone tags a release, not after.

---

## 7. Container Registries and GHCR

A **container registry** is a server that stores Docker images by name and tag, the same way GitHub stores code by commit — `docker push` uploads an image to one, `docker pull` downloads it, and any machine with network access to the registry (including your Chapter 88 cloud VM) can fetch the exact image you built in CI without ever rebuilding it from source. **GitHub Container Registry (GHCR)** is GitHub's own registry, and it is the natural choice here because it authenticates using the same GitHub identity and permissions your repository already has — no separate account to create, no separate credentials to manage.

Images in GHCR are named `ghcr.io/<owner>/<repo>`, so GoChain's published image ends up at `ghcr.io/you/gochain`. Anyone (or anything) with pull access can then run:

```bash
docker pull ghcr.io/you/gochain:v1.4.0
```

That single command replaces "SSH into the VM, git pull, rebuild the image from source" with "download the exact, already-tested image CI produced" — this is the mechanical core of what makes "push a tag" equivalent to "ship a release."

---

## 8. Writing release.yml — Tag, Build, Push

A separate workflow file handles releases, triggered only when you push a **Git tag** matching a version pattern like `v1.4.0`. Keeping this in its own file (rather than one giant `ci.yml` with conditionals everywhere) keeps each workflow's purpose obvious at a glance.

```yaml
# .github/workflows/release.yml
#
# Runs only when a tag matching v*.*.* is pushed — an ordinary commit to
# main never triggers this, only a deliberate, versioned release.
name: Release

on:
  push:
    tags:
      - "v*.*.*"

# Grants this workflow permission to push to GHCR under the repository's
# own identity — without this, the push step in the build job below
# would be rejected with a permissions error.
permissions:
  contents: read
  packages: write

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: "go.mod"
          cache: true
      # A release is exactly the moment you can least afford a skipped
      # test run, so the full suite runs again here even though it also
      # ran on the push to main that this tag presumably came from.
      - run: go test ./... -race -cover

  build-and-push:
    name: Build and push image
    runs-on: ubuntu-latest
    needs: test
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      # Extracts "v1.4.0" from the tag ref (refs/tags/v1.4.0) so it can be
      # used as the Docker tag below, without hardcoding a version anywhere.
      - name: Extract version
        id: version
        run: echo "value=${GITHUB_REF#refs/tags/}" >> "$GITHUB_OUTPUT"

      # Logs in to GHCR using the automatically-provided GITHUB_TOKEN —
      # scoped to this one workflow run, never a long-lived secret.
      - name: Log in to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      # Builds and pushes in one step, tagging the image both with the
      # exact version and with "latest," so a deploy script can always
      # pull the newest release without knowing its version number ahead
      # of time.
      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: |
            ghcr.io/you/gochain:${{ steps.version.outputs.value }}
            ghcr.io/you/gochain:latest
```

`secrets.GITHUB_TOKEN` is a special, automatically-generated credential GitHub Actions provides to every workflow run, scoped only to that run and only to the permissions declared in the `permissions:` block above — this is why no manual registry password ever needs to be created or stored for this step. `${{ github.actor }}` and `${{ steps.version.outputs.value }}` are GitHub Actions' expression syntax for referencing context values (who triggered the run) and outputs from earlier steps (the version string extracted two steps prior) — the same "reference a previous step's result" pattern you will see in almost every real-world workflow.

Cutting an actual release is now this, from your own terminal:

```bash
git tag v1.4.0
git push origin v1.4.0
```

Two commands, and a fresh, tested image lands in GHCR a few minutes later — no manual Docker build, no manual `docker push`, and critically, no way to publish a release that skipped the test suite.

---

## 9. From Registry to Running Node — the Deploy Step

Building and publishing the image is CI; getting the *running* node on your Chapter 88 VM to actually use it is the "CD" half. The simplest reliable approach is a deploy job that SSHes into the VM and re-pulls the Compose stack — no new infrastructure beyond what Chapter 88 already provisioned:

```yaml
  deploy:
    name: Deploy to testnet VM
    runs-on: ubuntu-latest
    needs: build-and-push
    steps:
      - name: Deploy over SSH
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.TESTNET_HOST }}
          username: ${{ secrets.TESTNET_SSH_USER }}
          key: ${{ secrets.TESTNET_SSH_KEY }}
          script: |
            cd /opt/gochain
            # Pulls the freshly-pushed :latest image referenced by
            # docker-compose.yml, then recreates only the containers
            # whose image actually changed — existing peers stay up.
            docker compose pull
            docker compose up -d
```

The `docker-compose.yml` on the VM references `ghcr.io/you/gochain:latest` (rather than building from local source), so `docker compose pull` fetches exactly the image `build-and-push` just published, and `docker compose up -d` recreates only the containers whose underlying image actually changed — a rolling update rather than tearing down the whole testnet at once. This is the entire "push a tag instead of manually redeploying" promise made good: `git tag` and `git push` are now the only two commands a maintainer ever needs to type to ship a new GoChain version to the live testnet.

---

## 10. Secrets: Keeping Credentials Out of the Repo

The deploy job needs an SSH host, username, and private key to reach the VM — none of which belong anywhere in the repository itself, since the repository (and every fork of it) is visible to anyone. GitHub's **repository secrets** (Settings → Secrets and variables → Actions) store these values encrypted, exposing them to workflows only as the `secrets.*` context, never printed in logs and never visible to a pull request from an outside contributor (GitHub deliberately withholds secrets from workflows triggered by forks, precisely to stop someone from opening a PR that exfiltrates your deploy key).

```
  Repository Settings → Secrets and variables → Actions
  --------------------------------------------------------
  TESTNET_HOST       = 203.0.113.42
  TESTNET_SSH_USER   = deploy
  TESTNET_SSH_KEY    = -----BEGIN OPENSSH PRIVATE KEY-----
                        ...
                        -----END OPENSSH PRIVATE KEY-----
```

Generate a dedicated deploy key (`ssh-keygen -t ed25519 -f deploy_key`) rather than reusing your personal key — a credential scoped to exactly one purpose is one you can revoke without disrupting anything else, and one whose blast radius, if it ever leaks, is limited to redeploying containers rather than full account access.

---

## 11. Watching a Pipeline Run

Every workflow run appears under your repository's **Actions** tab, showing each job, each step, and its output — exactly like watching the terminal output you would have gotten running these commands locally, except automatic and attached permanently to the commit or tag that triggered it. A red X on a pull request blocks it from looking "ready to merge" at a glance; a green check on a tag is your visible proof that the exact bits sitting in GHCR passed the full test suite before anyone deployed them.

```
  push commit  ─────────────►  CI: test  ──►  CI: build           (main)
                                    |
                                   pass
                                    |
                                    v
  push tag v1.4.0 ─────────────►  Release: test  ──►  build-and-push  ──►  deploy
                                                            |
                                                     ghcr.io/you/gochain:v1.4.0
                                                     ghcr.io/you/gochain:latest
```

---

## Summary

- CI runs the full test suite automatically on every push, so a broken change is caught in minutes, not after it reaches a live node.
- CD goes further: a tagged release automatically builds, tests, and publishes a Docker image with no manual steps.
- GitHub Actions workflows (`.github/workflows/*.yml`) are made of jobs, which run on isolated runners, made of steps (shell commands or reusable actions).
- `go test ./... -race -cover` runs every package's tests across all thirteen volumes in one command, with the race detector catching concurrency bugs GoChain's design invites by nature.
- `.github/workflows/ci.yml` runs on every push/PR: test, vet, format-check, then a Docker build-and-smoke-test.
- `.github/workflows/release.yml` runs only on `v*.*.*` tags: test again, build and push to GHCR tagged both with the version and `latest`, then deploy over SSH.
- `secrets.GITHUB_TOKEN` authenticates the registry push automatically; SSH credentials for the deploy step live in repository secrets, never in code.
- After this chapter, shipping a new GoChain version to the live testnet is exactly two commands: `git tag vX.Y.Z && git push origin vX.Y.Z`.

---

## Exercises

### Easy

1. Add `.github/workflows/ci.yml` from Section 5 to your own GoChain repository, push a small change, and paste the resulting Actions run's job list and pass/fail status.

2. Deliberately introduce a `gofmt` violation (misalign an indent by hand) and push it. Confirm the "Check formatting" step fails, then fix it and confirm a follow-up push passes.

3. Explain, in 3-4 sentences, why the build job in Section 6 uses `needs: test` instead of running independently and in parallel with the test job.

### Medium

4. Add a step to `ci.yml` that fails the build if `go vet` or the race detector finds anything in the `vm` package (Chapters 60-68) specifically, even if you decide to allow the rest of the suite to proceed on a warning — explain your reasoning for treating this package more strictly.

5. Extend `release.yml` so it only runs `build-and-push` if the tag was pushed from the `main` branch (reject tags created on unmerged feature branches). Describe the check you added and why it matters for a project where tags trigger real deployments.

6. The deploy job in Section 9 has no rollback step. Design (and, if you can, implement) a `rollback` job, triggered manually via `workflow_dispatch`, that re-deploys the *previous* tagged image instead of the latest one.

### Hard

7. Add a matrix build to `ci.yml` that runs the test job against two different Go versions (for example, the version in `go.mod` and the previous minor release), and explain what kind of bug this would catch that a single-version test run would miss.

8. The current `docker compose pull && docker compose up -d` deploy strategy briefly disrupts any node whose image changed. Research and describe (in 200-300 words) how you would extend this to a zero-downtime rolling deploy across multiple GoChain node containers, referencing Chapter 89's Kubernetes overview for how a `Deployment` resource handles this natively.

9. Design a `staging` environment: a second, smaller Compose stack that every push to `main` (not just tags) deploys to automatically, so changes are exercised on a real running network before anyone tags a production release. Sketch the workflow YAML changes needed, and explain what new secrets and VM this would require.
