# Parameterizer
[![Go Report Card](https://goreportcard.com/badge/github.com/kubernauts/parameterizer)](https://goreportcard.com/report/github.com/kubernauts/parameterizer)
[![godoc](https://godoc.org/github.com/kubernauts/parameterizer?status.svg)](https://godoc.org/github.com/kubernauts/parameterizer)

`Parameterizer` is a command line tool and Kubernetes operator for handling application lifecycle management, generically.

**Currently tested and works only with Kubernetes 1.9**

Just like `Ingress` allows to generically define how traffic is routed to Kubernetes services, with different backends (NGINX, HAProxy) providing the functionality, the `Parameterizer` resource defines a sequence of commands applied to an app definition input (directory or registry such as [Quay.io](https://quay.io/application/) or [Kubestack](https://www.kubestack.com/)), turning Kubernetes application definitions (e.g. expressed in [Helm templates](https://github.com/kubernetes/helm/blob/master/docs/chart_template_guide/functions_and_pipelines.md), [ksonnet](https://ksonnet.io/docs/concepts), [kapitan](https://github.com/deepmind/kapitan), etc.) along with user-defined parameters into a parameterized Kubernetes YAML manifest.

This parameterized YAML manifest defines the necessary deployments, services, etc. for the app and can, for example, be used in a `kubectl apply` command to create the resources or via Helms' [Tiller](https://docs.helm.sh/glossary/#tiller), [appr](https://github.com/app-registry/appr), or other installers/ALM tools.

The overall `Parameterizer` architecture is as follows:

![Parameterizer architecture](img/parameterizer-architecture.png)

## Install

For now we do not provide binaries so you'll need to have [Go](https://golang.org/dl/) installed to use the `krm` CLI tool. We've been testing it using `go1.9.2 darwin/amd64` and `go version go1.10 darwin/amd64`. 

To build `krm` from source, do the following:

```
$ go get github.com/kubernauts/parameterizer/cli/krm
```

Note that if your `$GOPATH/bin` is in your `$PATH` then now you can use `krm` from everywhere. If not, you can 1) do a `cd $GOPATH/src/github.com/kubernauts/parameterizer/cmd` followed by a `go build` and use it from this directory, or 2) run it using `$GOPATH/bin/krm`.

## Use

In general, the workflow would be something like:

1. the application author creates the `Parametrizer` manifest along with its package such as a `helm-chart.yaml`.
1. the operator can then take the `Parametrizer` resource and deploy the application with the installer or deploy manager she wants (e.g. `kubectl`).
1. in addition to above, the operator can create a new `Parametrizer` resource to chain additional transformations and/or or compose dependencies.

For example, if you have the following `Parameterizer` resource in a file `kubeless.yaml` ([source](examples/kubeless.yaml)):

```yaml
kind: Parameterizer
apiVersion: kubernetes.sh/v1alpha1
metadata:
  name: install-kubeless
spec:
  resources:
  - name: kubeless-jsonnet
    source:
      fetch:
        urls:
          - https://raw.githubusercontent.com/kubeless/kubeless/master/kubeless.jsonnet
          - https://raw.githubusercontent.com/kubeless/kubeless/master/kubeless-non-rbac.jsonnet
        dest: /resources
    volumeMounts:
      - name: kubeless-input
        mountPath: /resources

  volumes:
  - name: kubeless-input
    emptyDir: {medium: ""}

  transformations:
  - container:
      name: ksonnet
      image: quay.io/ant31/ksonnet
      command:
        - kubecfg
        - show
        - -o
        - yaml
        - /resources/kubeless.jsonnet
      volumeMounts:
        - name: kubeless-input
          mountPath: /resources
```

You can apply the parameters and install the app like so:

```
$ krm expand kubeless.yaml | kubectl apply -f -
```

## Test

For now, just the simple unit-level test in Go (that is, no integration tests yet):

```
$ cd $GOPATH/src/github.com/kubernauts/parameterizer/pkg/parameterizer
$ go test
PASS
ok      github.com/kubernauts/parameterizer/pkg/parameterizer   0.007s
```
