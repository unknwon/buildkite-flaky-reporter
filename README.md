# buildkite-flaky-reporter

HTTP server for reporting flaky tests on Buildkite. In addition, the reporter could also be used to notify when the job is canceled.

## Configuration

By default, it is located in current directory as `app.ini`. Or use `-config` flag to specify another one.

```
addr = localhost:4222

[buildkite]
token = {REDACTED}
org_slug = sourcegraph
pipeline_slug = sourcegraph
build_branch = master
job_name = :chromium:
failures_threshold = 3
report_cancel = true

[slack]
url = {REDACTED}
cancel_notify_url = {REDACTED}
```

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.
