[](https://img.shields.io/badge/version-development-red.svg)
# snap-plugin-collector-portcount

This plugin counts open ports of responding hosts based on [redforks portscan](https://github.com/redforks/portscan)

It's used in the [Snap framework](http://github.com:intelsdi-x/snap).

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
3. [License](#license-and-authors)
4. [Acknowledgements](#acknowledgements)

## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/)  - needed only for building. See also [How to install Go language](http://ask.xmodulo.com/install-go-language-linux.html)

### Operating systems
All OSs currently supported by Snap:
* Linux/amd64

### Installation
#### To build the plugin binary:
```
$ go get -u github.com/IrekRomaniuk/snap-plugin-collector-portcount
```
### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Load the plugin and create a task, see example in [Examples](https://github.com/IrekRomaniuk/snap-plugin-collector-portcount/tree/master/examples).

## Documentation

### Collected Metrics

This plugin has the ability to gather the following metric:

Namespace | Description
----------|-----------------------
/niuk/portcount/* | total number of hosts responding on given port (per config in task manifest)


### Example
Example running portcount collector and writing data to an Influx database.

Load portcount plugin
```
$ snaptel plugin load $GOPATH/bin/snap-plugin-collector-portcount
```
List available plugins
```
$ snaptel plugin list
NAME                     VERSION         TYPE            SIGNED          STATUS          LOADED TIME
portscan                 1               collector       false           loaded          Sun, 11 Dec 2016 13:02:43 EST
```
See available metrics for your system
```
$ snaptel metric list
```

Create a task manifest file (see below) and put full path to the [file](https://github.com/IrekRomaniuk/snap-plugin-collector-portcount/blob/master/examples/pinglist.txt) listing IP addresses:
```yaml
deadline: "15s"
version: 1
schedule:
  type: "simple"
  interval: "30s"
max-failures: 10
workflow:
  collect:
    metrics:
      /niuk/portcount/*: {}
    config:
      /niuk/portcount:
        target: "/global/path/iplist.txt"
        port: "53"
```
Load influxdb plugin for publishing:
```
$ snaptel plugin load snap-plugin-publisher-influxdb
```

Create a task:
```
$ snaptel task create -t task-portcount.yml -n portcount-981
Using task manifest to create task
Task created
ID: 30d8c7d1-4cec-4eec-a787-3eaa2a088436
Name: portcount-981
```

List running tasks:
```
$ snaptel task list
ID                                       NAME                                            STATE           HIT     MISS    FAIL    CREATED                 LAST FAILURE          
30d8c7d1-4cec-4eec-a787-3eaa2a088436     portcount-981                                   Running         33      0       0       1:04PM 12-11-2016   
```
Watch the task
```
$ snaptel task watch 30d8c7d1-4cec-4eec-a787-3eaa2a088436
Watching Task (30d8c7d1-4cec-4eec-a787-3eaa2a088436):
NAMESPACE                DATA    TIMESTAMP
/niuk/portcount/981       305     2016-12-11 13:19:47.785976086 -0500 EST
```
Watch metrics in real-time using [Snap plugin for Grafana] (https://blog.raintank.io/using-grafana-with-intels-snap-for-ad-hoc-metric-exploration/)

## License
This plugin is Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@IrekRomaniuk](https://github.com/IrekRomaniuk/)



