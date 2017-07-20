# Ksonnet

Yet another templatizing configurator for Kubernetes resources. Based on Jsonnet,
uses ksonnet-lib extension for writing configuration jsonnet template, which
after that rendered with jsonnet and output Kubernetes template for further usage
with `kubectl create`.

Ksonnet home page: http://ksonnet.heptio.com/

## Difference from Helm

Ksonnet has next differences from Helm package tool:

* Ksonnet allows not to write all template, as it happened in Helm, it uses ksonnet-lib
  functions for this.

* Ksonnet has convenient tool for rendering named kubecfg (https://github.com/ksonnet/kubecfg).
  By default ksonnet templates rendered by jsonnet.

* Ksonnet templates rendered before deploying kubernetes resources, helm templates
  rendered together with deploying template resources. But kubecfg tool has
  required CRUD cli, like Helm client.

* To change variable values in Ksonnet, need to open jsonnet template file and
  manually add changes. In Helm there are values, which could be changed without
  changing templates for deploying different cases.

## Ksonnet disadvantages

* Ksonnet has no adequate documentation - need to learn ksonnet methods from source
  code.

* If ksonnet doesn't contain required function - need to write raw kubernetes templates
  as a variable.

## How to install

Following steps allows to install Ksonnet:

1. First, install jsonnet:

   ```console
   git clone https://github.com/google/jsonnet; cd jsonnet
   make
   mv jsonnet /usr/bin/jsonnet
   ```

2. Then install ksonnet:

   * Fork or clone this repository, using a command such as:
     ```console
     git clone https://github.com/ksonnet/ksonnet-lib
     ```
   * Then add the appropriate import statements for the library to your Jsonnet code:
     ```console
     local k = import "ksonnet.beta.1/k.libsonnet";
     ```

## Examples

kafka/ and zookeeper/ directories demonstrates ksonnet templates usage. Please,
address the examples.

## Potential

Ksonnet can be templatized with different values with different templatizators
like jinja2 etc.

## Conclusion

Ksonnet is convenient language only if you need to write small and fast kubernetes resource
template. Otherwise, ksonnet template take much more place and time to write over against
Helm.
