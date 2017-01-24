<!--
http://www.apache.org/licenses/LICENSE-2.0.txt


    Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# snap-plugin-publisher-elasticsearch

This plugin supports publishing metrics to Elasticsearch endpoint

It is used in the [Snap framework] (http://github.com/intelsdi-x/snap).


1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

### System Requirements

* [golang 1.6+](https://golang.org/dl/)  (needed only for building)
* Elasticsearch instance to publish data (docker run elasticsearch)

### Installation

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-publisher-elasticsearch  
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-publisher-elasticsearch.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```

This builds the plugin in `./build/`

### Configuration and Usage

* Set up the [Snap framework](https://github.com/intesdi-x/snap/blob/master/README.md#getting-started)

## Documentation

The plugin expects you to provide following parameters:
 - `protocol` http, https, tls
 - `address` address of elasticsearch endpoint
 - `port` elasticsearch port
 - `index` index for publishing data
 - `type` mapping type

You can also set following options if needed:
 - `index_timestamp` default is set to `true` (boolean). Set to false if not want to extend index prefix with timestamp
 - `publish_fields` default is set to `Namespace|Config|Data|Description|Timestamp|Unit|Version` (string). Remove field from metric structure which shouldn't be published.

Plugin publishes dynamic elements from metric namespace.

### Examples

You can find exemplary task manifests for using this plugin in folder [examples](examples/tasks).

Example running [logs collector plugin](https://github.com/intelsdi-x/snap-plugin-collector-logs), [logs-regexp processor](https://github.com/intelsdi-x/snap-plugin-processor-logs-regexp),
[logs-openstack processor](https://github.com/intelsdi-x/snap-plugin-processor-logs-openstack), and writing data to elasticsearch.

Prepare docker with elasticsearch
```
$ docker pull elasticsearch
$ docker run --name elasticsearch -p 9200:9200 -p 9300:9300 elasticsearch
```

In one terminal window, open the Snap daemon:
```
$ snapteld -t 0 -l 1 -m 999
```

The option "-l 1" it is for setting the debugging log level, "-t 0" is for disabling plugin signing, and "-m 999" is for the maximum number of instances of a loaded plugin to run.

In another terminal window:

Download and load collector, processor and publisher plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-logs/latest/linux/x86_64/snap-plugin-collector-logs
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-processor-logs-regexp/latest/linux/x86_64/snap-plugin-processor-logs-regexp
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-processor-logs-openstack/latest/linux/x86_64/snap-plugin-processor-logs-openstack
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-elasticsearch/latest/linux/x86_64/snap-plugin-publisher-elasticsearch
$ snaptel plugin load snap-plugin-collector-logs
$ snaptel plugin load snap-plugin-processor-logs-regexp
$ snaptel plugin load snap-plugin-processor-logs-openstack
$ snaptel plugin load snap-plugin-publisher-elasticsearch
```

Find logs to collect or use exemplary logs available in [logs](examples/logs).

Create task manifests, see examples:
- to use snap-plugin-collector-logs, snap-plugin-processor-logs-regexp, snap-plugin-publisher-elasticsearch - [logs-regexp-elasticsearch.json](examples/tasks/logs-regexp-elasticsearch.json)
- to use snap-plugin-collector-logs, snap-plugin-processor-logs-openstack, snap-plugin-publisher-elasticsearch - [logs-openstack-elasticsearch.json](examples/tasks/logs-openstack-elasticsearch.json)

If you are using exemplary logs you need to modify `"log_dir"` in task manifests and set appropriate paths to directories with logs.
If you want to collect your own logs you need to adjust config sections in task manifests to your environment.

Create tasks:
```
$ snaptel task create -t logs-regexp-elasticsearch.json
$ snaptel task create -t logs-openstack-elasticsearch.json
```

See running tasks:
```
$ snaptel task list
```

Check indexes in elasticsearch:
```
$ curl 'http://localhost:9200/_cat/indices?v'
```

See data saved in elasticsearch:
```
$ curl -XGET 'http://localhost:9200/snap*/message/_search?pretty=true&size=9999'
```

Exemplary data saved in elasticsearch has the form which is shown in [outputs](examples/outputs).

### Roadmap

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-publisher-elasticsearch/issues) 
and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-publisher-elasticsearch/pulls).

## Community Support
This repository is one of **many** plugins in the **Snap**, a powerful telemetry agent framework. See the full project at 
http://github.com/intelsdi-x/snap. To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).


## Contributing
We love contributions! :heart_eyes:

There is more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).


## Acknowledgements

* Author: [Patryk Matyjasek](https://github.com/PatrykMatyjasek)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
