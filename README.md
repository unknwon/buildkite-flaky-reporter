# buildkite-flaky-reporter
HTTP server for reporting flaky tests on Buildkite

## Configuration

```
addr = localhost:4222

[buildkite]
token = {REDACTED}
org_slug = sourcegraph
pipeline_slug = sourcegraph
build_branch = master
job_name = :chromium:
failures_threshold = 3

[slack]
url = {REDACTED}
```
