[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/DoctorOgg/sensu-check-statuspage-betterstack)
![goreleaser](https://github.com/DoctorOgg/sensu-check-statuspage-betterstack/workflows/goreleaser/badge.svg)

# sensu-check-statuspage-betterstack

## Table of Contents
- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)

## Overview

The sensu-check-statuspage-betterstack is a [Sensu Check][6] that will check a status page hosted by betterstack. 


## Files

* sensu-check-statuspage-betterstack

## Usage examples

```bash
sensu-check-statuspage-betterstack -u https://status.bunnycdn.com

bunny.net: Incidents: 0, Updated at: 2021-06-23 07:25:05.383 +0000 UTC
```

```bash
sensu-check-statuspage -u https://status.ucdavis.edu

UC Davis: Incidents: 2, Updated at: 2021-07-10 08:30:33.455 -0700 PDT
MAJOR: Voice Service Degradation: Jabber clients for VOIP (https://stspg.io/mk26h5mgpxyf) IDENTIFIED
MINOR: ServiceNow Degradation (https://stspg.io/csf4p3923xwr) INVESTIGATING
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add DoctorOgg/sensu-check-statuspage
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/DoctorOgg/sensu-check-statuspage-betterstack].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-check-statuspage
  namespace: default
spec:
  command: sensu-check-statuspage-betterstack --url https://status.example.com
  subscriptions:
  - system
  runtime_assets:
  - DoctorOgg/sensu-check-statuspage-betterstack

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-check-statuspage repository:

```
go build
```


[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu-community/check-plugin-template/blob/master/.github/workflows/release.yml
[5]: https://github.com/sensu-community/check-plugin-template/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[7]: https://github.com/sensu-community/check-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu-community/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
