# Parameterizer

`Parameterizer` is a Kubernetes operator for handling application lifecycle management, generically. Just like `Ingress` allows to generically define how traffic is routed to Kubernetes services, with different backends (NGINX, HAProxy) providing the functionality, the `Parameterizer` resource allows you to define a sequence of commands applied to an input (directory or registry such as Quay or Kubestack), turning Kubernetes application definitions (e.g. expressed in Helm templates, ksonnet, kapitan, jinja2, etc.) along with parameters provided by the user into a set of parameterized Kubernetes YAML manifests, defining deployments, services, etc.—this output can, for example, be used in a `kubectl apply` command to create the resources or via installers such as Tiller or CoreOS ALM.

## Install

## Use
