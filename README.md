# Otterize Helm Charts

![build](https://github.com/otterize/helm-charts/actions/workflows/build.yaml/badge.svg)
[![community](https://img.shields.io/badge/slack-Otterize_Slack-orange.svg?logo=slack)](https://join.slack.com/t/otterizeworkspace/shared_invite/zt-1fnbnl1lf-ub6wler4QrW6ZzIn2U9x1A)

This repository contains the official Helm charts for configuring and installing Otterize OSS components on Kubernetes: the intents operator, the network mapper, and the SPIRE integration operator. These charts support multiple use cases depending on the values provided.

For full documentation on installing and configuring Otterize with these Helm charts, as well as all the ways you can use Otterize with Kubernetes, please see the
[docs site](https://docs.otterize.com/).

## Prerequisites

To use the charts here, you'll need [Kubernetes](https://kubernetes.io/docs/home/) and [Helm](https://helm.sh/docs/intro/quickstart/).

## Usage

To use the latest version of these charts, first add the Otterize helm repository:

```console
$ helm repo add otterize https://helm.otterize.com
```
You should see:
```console
"otterize" has been added to your repositories
````

Then you can install any of them into your cluster, e.g. to install the otterize-kubernetes chart into a namespace called "otterize-system", use:
```console
$ helm install otterize otterize/otterize-kubernetes -n otterize-system
```

The README for each chart describes the various configuration options it supports. 
These are also documented in the [docs site](https://docs.otterize.com/) along with more detailed installation instructions.
