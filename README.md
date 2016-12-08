# Snap publisher plugin - Circonus

This plugin publishes metrics to circonus.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
1. [Documentation](#documentation)
  * [Task Manifest Config](#task-manifest-config)
  * [Examples](#examples)
1. [License](#license)

## Getting Started

### System Requirements
* The Snap daemon is running
* A [Circonus](https://circonus.com/) account or Circonus Inside installation reachable by the plugin is required for successful publishing of metrics.

### Installation

#### Download Circonus plugin binary:
You can get the pre-built binaries for your OS and architecture at plugin's [Github Releases](https://github.com/circonus-labs/snap-plugin-publisher-circonus/releases) page.

#### To build the plugin binary:

##### Build requirements
* [`govendor`](https://github.com/kardianos/govendor)

```
mkdir -p "${GOPATH}/src/github.com/circonus-labs"
cd "${GOPATH}/src/github.com/circonus-labs"
git clone https://github.com/circonus-labs/snap-plugin-publisher-circonus
cd snap-plugin-publisher-circonus
make
```

This builds the plugin in `./build`

### Configuration and Usage

* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

## Documentation

### Task Manifest Config

A Task Manifest that includes publishing to Circonus will require configuration data in order for the plugin to establish a connection. Config options are detailed in [OPTIONS.md](OPTIONS.md).

#### Metric types

Circonus supports three main metric **value** types: numeric, histogram, and text. Which of these to use can be controlled in the `tags` task element, `"circonus_type": "(numeric|histogram|text)"`. The default is to send metric values as numeric. The example task manifests illustrate setting explicit numeric and histogram tags for the metrics being collected.

e.g. (the collected metrics are histograms see [`example/tasks/tasks-histogram.json`](example/tasks/tasks-histogram.json))

```json
"tags": {
    "/intel/procfs": {
        "circonus_type": "histogram"
    }
}
```

### Examples

See the [example](example/) directory for a complete, working example using a Vagrant VM.

## License

See [License](LICENSE)
