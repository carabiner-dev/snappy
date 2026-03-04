# Snappy
### Grab attestable snapshots of API responses

Snappy is a simple tool that grabs snapshots of API response data to
be able to fix it in a signed attestation. It supports both **GitHub**
and **GitLab** APIs.

The API call is constructed according to a spec file. To generate a snapshot
you simply call snappy and pass it a spec:

```bash
snappy snap spec.yaml
```

## Supported Platforms

Snappy supports the following platforms:

| Platform | Auth Environment Variable | Host Override |
|----------|--------------------------|---------------|
| GitHub   | `GITHUB_TOKEN`           | `GITHUB_HOST` |
| GitLab   | `GITLAB_TOKEN`           | `GITLAB_HOST` |

The platform is auto-detected from the spec's endpoint, or can be
explicitly set with the `--platform` flag:

```bash
snappy snap --platform gitlab spec.yaml
```

Both platforms support self-hosted instances through their respective
host environment variables.

## Snapshot Spec

The API snapshots are controlled through a spec file written in YAML.
The spec describes some basic metadata about the object being queried,
and the API endpoint to hit.

In order to avoid leaking sensitive data, only data specifically configured
in the spec will be written in the resulting snapshot.

### GitHub Example

```yaml
---
id: https://github.com/${ORG}/${REPO}
type: http://github.com/carabiner-dev/snappy/specs/repo.yaml
endpoint: repos/${ORG}/${REPO}
method: GET
mask:
  - id
  - node_id
  - name
  - full_name
  - html_url
  - git_url
  - license
  - visibility
  - security_and_analysis
```

### GitLab Example

```yaml
---
id: ${HOST}/${GROUP}/${PROJECT}
type: http://github.com/carabiner-dev/snappy/specs/gitlab/project.yaml
endpoint: projects/${GROUP}%2F${PROJECT}
method: GET
mask:
  - id
  - name
  - visibility
  - default_branch
  - merge_method
  - approvals_before_merge
```

From the resulting data structure, only those fields named in the field mask
will be transferred to the resulting snapshot.

### Built-in Specs

Snappy ships with built-in specs that can be referenced with the `builtin:` prefix:

**GitHub:** `repo`, `org`, `commit`, `branch-rules`, `mfa`

**GitLab:** `project`, `branch-protection`

```bash
snappy snap --var ORG=myorg --var REPO=myrepo builtin:github/repo
```

## Variable Support

As seen in the examples above, the snapshot spec can have variables that
can be substituted when invoking snappy:

```yaml
id: https://github.com/${ORG}/${REPO}
endpoint: repos/${ORG}/${REPO}
```

To define a value for the `$ORG` and `$REPO` variables, simply pass them
in the command line invocation:

```bash
snappy snap --var ORG=myorg --var REPO=myrepo spec.yaml
```

## Install

### Pre-built Binaries

Download a pre-built binary for your platform from the
[GitHub Releases](https://github.com/carabiner-dev/snappy/releases) page.

Binaries are available for Linux (amd64, arm64), macOS (arm64), and
Windows (amd64).

### Using Go

If you have Go installed, you can install snappy with:

```bash
go install github.com/carabiner-dev/snappy@latest
```

## Patches Always Welcome

Snappy is released under the Apache 2.0 license by Carabiner Systems, Inc. Feel
free to contribute patches or open an issue if you find a problem. We love feedback!
