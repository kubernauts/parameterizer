# Parameterizer
[![Go Report Card](https://goreportcard.com/badge/github.com/kubernauts/parameterizer)](https://goreportcard.com/report/github.com/kubernauts/parameterizer)
[![godoc](https://godoc.org/github.com/kubernauts/parameterizer?status.svg)](https://godoc.org/github.com/kubernauts/parameterizer)

`Parameterizer` is a command line tool and Kubernetes operator for handling application lifecycle management, generically.

Just like `Ingress` allows to generically define how traffic is routed to Kubernetes services, with different backends (NGINX, HAProxy) providing the functionality, the `Parameterizer` resource defines a sequence of commands applied to an app definition input (directory or registry such as [Quay.io](https://quay.io/application/) or [Kubestack](https://www.kubestack.com/)), turning Kubernetes application definitions (e.g. expressed in [Helm templates](https://github.com/kubernetes/helm/blob/master/docs/chart_template_guide/functions_and_pipelines.md), [ksonnet](https://ksonnet.io/docs/concepts), [kapitan](https://github.com/deepmind/kapitan), etc.) along with user-defined parameters into a parameterized Kubernetes YAML manifest.

This parameterized YAML manifest defines the necessary deployments, services, etc. for the app and can, for example, be used in a `kubectl apply` command to create the resources or via Helms' [Tiller](https://docs.helm.sh/glossary/#tiller), [appr](https://github.com/app-registry/appr), or other installers/ALM tools.

The overall `Parameterizer` architecture is as follows:

![Parameterizer architecture](img/parameterizer-architecture.png)

## Install

For now we do not provide binaries so you'll need to have [Go](https://golang.org/dl/) installed to use the `krm` CLI tool. We've been testing it using `go1.9.2 darwin/amd64`. To build `krm` from source, do the following:

```
$ go get github.com/kubernauts/parameterizer/cmd/krm
```

Note that if your `$GOPATH/bin` is in your `$PATH` then now you can use `krm` from everywhere. If not, you can 1) do a `cd $GOPATH/src/github.com/kubernauts/parameterizer/cmd` followed by a `go build` and use it from this directory, or 2) run it using `$GOPATH/bin/krm`.

## Use

For example, if you have the following `Parameterizer` resource in a file `install-ghost-with-helm.yaml` ([source](test/install-ghost-with-helm.yaml)):

```
kind: Parameterizer
apiVersion: kubernetes.sh/v1alpha1
metadata:
  name: install-ghost
spec:
  resources:
  - name: helm-chart
    source:
      urls:
        - https://github.com/kubernetes/charts/tree/master/stable/ghost
    volume:
    - name: helm-chart
      hostPath:
        path: /tmp/
  - name: local-kinflate
    source:
      hostPath: ./resources/
    volume:
    - name: kinflate
  userInputs:
  - name: helm-user-values
    source:
       hostPath:
         path: ./values/prod
    volume:
    - name: helm-user-values
  volumes:
  - name: helm-output
    emptyDir: {medium: ""}
  apply:
  - image: lachlanevenson/k8s-helm:v2.7.2
    commands:
    -  helm template charts -f /helm-values/value.yaml -o /output/ghost-resources.yaml
    volumeMounts:
    - name: helm-output
      mountPath: /output
    - name: helm-chart
      mountPath: /charts
    - name: helm-user-values
      mountPath: /helm-values
  - image: ant31/kinflate
    commands:
    - bash -c 'cp /output/*.yaml" /kinflate/resources/all-resource.yaml" && kinflate inflate -f /kinflate'
    volumeMounts:
    - name: helm-output
      mountPath: /output
    - name: local-kinflate
      mountPath: /kinflate
```

You can apply the parameters and install the app like so:

```
$ krm expand install-ghost-with-helm.yaml | kubectl apply -f -
```
