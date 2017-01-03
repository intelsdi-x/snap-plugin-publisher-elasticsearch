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

You can find example task manifest for using this plugin in folder [examples](https://github.com/intelsdi-x/snap-plugin-publisher-elasticsearch/examples)

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
