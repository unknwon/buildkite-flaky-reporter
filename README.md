# buildkite-flaky-reporter

HTTP server for reporting flaky tests and cancellation on Buildkite.

## Preview

![image](https://user-images.githubusercontent.com/2946214/74906161-20be8c00-53eb-11ea-8a97-62e3c7da4198.png)
![image](https://user-images.githubusercontent.com/2946214/74906186-316f0200-53eb-11ea-8bf7-fd346e29de4e.png)

## Configuration

By default, it is located in current directory as `app.ini`. Or use `-config` flag to specify another one.

```ini
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
