## Example demonstrating the Circonus publisher plugin

Simple [Vagrant](https://vagrantup.com/) based example which will:

1. Start a CentOS 7 VM
2. Install [Snap telemetry framework](https://github.com/intelsdi-x/snap)
3. Download and install [Go v1.7.4](https://golang.org)
4. Download and install several system level [snap collector plugins](https://github.com/intelsdi-x/snap/blob/master/docs/PLUGIN_CATALOG.md) (e.g. cpu, memory, disk, fs, interface)
5. Clone [this repository](https://github.com/circonus-labs/snap-plugin-publisher-circonus)
6. Build and install the `snap-plugin-publisher-circonus` plugin binary
7. Configure two sets of tasks and (re)start the `snap-telemetry` service

### Setup

Set an environment variable with a valid [Circonus API Token](https://login.circonus.com/user/tokens).
```sh
export CIRCONUS_API_TOKEN="..."
```

#### Start/Run

```sh
vagrant up
```

##### Manage [running] VM
```sh
# access
vagrant ssh
# stop
vagrant halt
# destroy
vagrant destroy
```

### Customize

Edit the `Vagrantfile` and/or `tasks/tasks-*.json` files to further customize if needed. Add plugins, test various Circonus-specific plugin [options](../OPTIONS.md), etc.
