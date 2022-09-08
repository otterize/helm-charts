# Otterize Helm Charts

This repository contains the official Helm charts for configuring and installing Otterize open source components on Kubernetes: the Intents Operator, Network Mapper and SPIRE Integration Operator. These charts support multiple use cases depending on the values provided.

For full documentation on this Helm chart along with all the ways you can use Otterize with Kubernetes, please see the
[Otterize documentation](https://docs.otterize.com/).

## Prerequisites

To use the charts here, you need [Helm](https://helm.sh/). Setting up Kubernetes and Helm is outside the scope of this README. Please refer to the Kubernetes and Helm documentation.

## Usage

To install the latest version of these charts, add the Otterize helm repository
and run `helm install`:

```console
$ helm repo add otterize https://helm.otterize.com
"otterize" has been added to your repositories

$ helm install otterize -n otterize-deploy otterize/otterize-kubernetes # or another chart
```

Please see the many options supported in the README for each chart. These are also
fully documented directly on the [Otterize website](https://docs.otterize.com/) along with more detailed installation instructions.
