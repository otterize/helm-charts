# Otterize Helm Charts

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

Then you can install any of them into your cluster, e.g. to install the otterize-kubernetes component into a namespace called "otterize-deploy", use:
```console
$ helm install otterize otterize/otterize-kubernetes -n otterize-deploy
```

The README for each chart describes the various configuration options it supports. 
These are also documented in the [docs site](https://docs.otterize.com/) along with more detailed installation instructions.
