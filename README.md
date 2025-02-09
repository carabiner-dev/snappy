# Snappy
### Grab attestable snapshots of API responses

Snappy is a simple tool that grabs snapshots of API response data to
be able to fix it in a signed attestation.

The API call is contructed according to a spec file, to generate a snapshot
you simply call snappy and pass it a spec:

```
snappy snap spec.yaml
```

## Snapshot Spec 

The API snapshots are controlled through a spec file written in yaml.
The spec describes some basic metadata about the object being queried,
and the API endpoint to hit.

In order to avoid leaking snesitive data, only data specifically configured
in the spec will be written in the resulting snapshot.

An example spec YAML file looks like this:

```yaml
---
id: https://github.com/{$ORG}/{$REPO}
type: http://github.com/carabiner-dev/snapi/specs/repo.yaml
endpoint: repos/{$ORG}/{$REPO}
method:  GET
mask:
  - id
  - node_id
  - name
  - full_name
  - html_url
  - git_url
  - license
  - has_discussions
  - has_issues
  - homepage
  - fork
  - visibility
  - security_and_analysis
  - web_commit_signoff_required
```

From the resulting data structure, only those fields named in the field mask
will be transferred to the resulting snapshot.

## Variable Support

As seen in the example above, the snapshot spec can have variables that
can be substituted when invoking snappy:

```yaml
id: https://github.com/{$ORG}/{$REPO}
type: http://github.com/carabiner-dev/snapi/specs/repo.yaml
endpoint: repos/{$ORG}/{$REPO}
```

To define a value for the `$ORG` and `$REPO` variables, simply pass them
in the command line invocation:

```
snappy snap --var ORG=myorg --var REPO=myrepo spec.yaml
```

## Install

TBD

