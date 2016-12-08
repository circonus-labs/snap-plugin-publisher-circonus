## Circonus publisher for the [Snap telemetry framework](https://github.com/intelsdi-x/snap)

### Example 'publish' sections of task workflow

Minimum required options:

```json
"publish": [
    {
        "plugin_name": "circonus",
        "config": {
            "api_token": "..."
        }
    }
]
```

All available options (with defaults, exception is `check_instance_id` which is generated based on the result of `os.Hostname()`):

```json
"publish": [
    {
        "plugin_name": "circonus",
        "config": {
            "interval": "10s",
            "reset_counters": "true",
            "reset_gauges": "true",
            "reset_histograms": "true",
            "reset_text": "true",
            "log_level": "warn",
            "api_token": "",
            "api_app": "snap-cgm",
            "api_url": "https://api.circonus.com/v2",
            "check_id": "",
            "check_submission_url": "",
            "check_instance_id": "",
            "check_target_host": "",
            "check_display_name": "",
            "check_search_tag": "service:snap-telemetry",
            "check_secret": "",
            "check_tags": "",
            "check_max_url_age": "",
            "check_force_metric_activation": "false",
            "broker_id": "",
            "broker_select_tag": "",
            "broker_max_response_time": "500ms"
        }
    }
]
```

## Options
| Option | Default | Description |
| ------ | ------- | ----------- |
| General ||
| `interval` | "10s" | Interval at which metrics are flushed and sent to Circonus. |
| `reset_counters` | "true" | Reset counter metrics after each submission. |
| `reset_gauges` | "true" | Reset gauge metrics after each submission. |
| `reset_histograms` | "true" | Reset histogram metrics after each submission. |
| `reset_text` | "true" | Reset text metrics after each submission. |
| `log_level` | "warn" | Controls level of logging messages emitted. |
|API||
| `api_token` | "" | [Circonus API Token key](https://login.circonus.com/user/tokens) |
| `api_app` | "snap-cgm" | App associated with API token |
| `api_url` | "https://api.circonus.com/v2" | Circonus API URL |
|Check||
| `check_id` | "" | ID of previously created check. (*Note: **check id** not **check bundle id**.*) |
| `check_submission_url` | "" | Submission URL of previously created check. |
| `check_instance_id` | hostname:snap-telemetry | An identifier for the 'group of metrics emitted by this process or service'. Hostname is whatever `os.Hostname()` returns. |
| `check_target_host` | check_instance_id | Explicit setting of `check.target`. |
| `check_display_name` | check_instance_id | Custom `check.display_name`. Shows in UI check list. |
| `check_search_tag` | "service:snap-telemetry" | Specific tag used to search for an existing check when neither SubmissionURL nor ID are provided. |
| `check_tags` | "" | List (comma separated) of tags to add to check when it is being created. The SearchTag will be added to the list. |
| `check_secret` | random generated | A secret to use for when creating an httptrap check. |
| `check_max_url_age` | "5m" | Maximum amount of time to retry a [failing] submission URL before refreshing it. |
| `check_force_metric_activation` | "false" | If a metric has been disabled via the UI the default behavior is to *not* re-activate the metric; this setting overrides the behavior and will re-activate the metric when it is encountered. |
|Broker||
| `broker_id` | "" | ID of a specific broker to use when creating a check. Default is to use a random enterprise broker or the public Circonus default broker. |
| `broker_select_tag` | "" | Used to select a broker with the same tag(s). If more than one broker has the tag(s), one will be selected randomly from the resulting list. (e.g. could be used to select one from a list of brokers serving a specific colo/region. "dc:sfo", "loc:nyc,dc:nyc01", "zone:us-west") |
| `broker_max_response_time` | "500ms" | Maximum amount time to wait for a broker connection test to be considered valid. (if latency is > the broker will be considered invalid and not available for selection.) |

## Notes:

* All options are *strings*.
* At a minimum, one of either `api_token` or `check_submission_url` is **required** for correct functionality.
* Check management can be disabled by providing a `check_submission_url` without an `api_token`. Note: the supplied URL needs to be http or the broker needs to be running with a cert that can be verified. Otherwise, the `api_token` will be required to retrieve the correct CA certificate to validate the broker's cert for the SSL connection.
* A note on `instance_id`, the instance id is used to consistently identify a check. The display name can be changed in the UI. The hostname may be ephemeral. This is more applicable to applications developed using [cgm](https://github.com/circonus-labs/circonus-gometrics) rather than Snap. Where metric continuity is important across execution boundaries, the instance id is used to locate existing checks. Since the check.target is never actually used by an httptrap check it is more decorative than functional, a valid FQDN is not required for an httptrap check.target. But, using instance id as the target can pollute the Host list in the UI with application specific entries.
