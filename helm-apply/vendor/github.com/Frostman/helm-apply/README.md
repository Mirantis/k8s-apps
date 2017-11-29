# helm-apply

## DSL overview

```yaml
repos:
  <chart-repo-name>: <chart-repo-url>

clusters:
  <cluster-name>:
    context: <context-name> # Optional. Current context by default.
    path: <path-to-kube-config> # Optional. Kube config path will be discovered by default.
    namespace: <ns> # Optional. Namespace to install releases into. Current context namespace by default.
    tillerNamespace: <ns> # Namespace of tiller. Will try to find it in `kube-system` if not specifyed.
                          # If tiller was not found, it will be installed in `tillerNamespace` (or `kube-system` namespace if `tillerNamespace` is undefined)
    externalIP: <IP> #Optional. Ip that will be used for cross-cluster communication if service type is NodePort.

releases:
  <release-name>:
    chart: <chart-repo-name>/<chart-name>:<chart-version> # Required. Chart version is optional - latest will be used by default.
    cluster: <cluster-name> # Required.
    namespace: <ns> # Namespace to install release into. Uses cluster `namespace` by default.
    wait: true # Optional. Wait untill release becomes ready. False by default.
    dependencies: # Optional. Please refer to `Dependencies` section.
      <dep-chart-name>: <dep-release-name>
    parameters: # Optional. Overrides values.yaml of the chart.
      <some-param>: <some-value>
```

## Chart requirements

There are no any specific requirements for charts in general. But if you want
to use `dependencies` feature, you have to modify your char in the way
it described in [Dependencies][].

## Config examples

[Minimalistic config with one cluster and two releases.](./examples/config.yaml)

[Twitter demo sample config with two clusters.](./examples/twitter-demo.yaml)

## Dependencies

Dependencies between charts are currently based on:
 * NOTES output of release to get internal or external addresses of
   exposed services
 * Helm native Chart-requirements model
 * Mechanisms to disable child Charts deployment and propagate addresses
   through the Values instead.

There are two example charts that support dependencies:
 * [Tweeviz API][https://github.com/Mirantis/k8s-apps/tree/master/charts/tweeviz_api]
 * [Tweeviz UI][https://github.com/Mirantis/k8s-apps/tree/master/charts/tweeviz_ui]

To expose services provided by chart, NOTES.txt should have the following:

```
Internal URL:
  <name>: <url>
```
If those services need to be publicly accessible, the following should
be added as well:
```
External services:
  <name>: <k8s-service-name>:<k8s-port or k8s-port-name>
```

Internal urls will be used to connect charts of the same K8s cluster. To
connect charts deployed on different K8s clusters, url will be formed
based on type of the service provided in `External services`.

In our example we have the following in NOTES.tx:

For tweeviz_ui:
```
Internal URL:
  tweeviz_ui: {{ template "tweeviz.ui.fullname" . }}:{{ .Values.port }}

External services:
  tweeviz_ui: {{ template "tweeviz.ui.fullname" . }}:{{ .Values.port }}
```

For tweeviz_api:
```
Internal URL:
  tweeviz_api: {{ template "tweeviz.api.fullname" . }}:{{ .Values.port }}

External services:
  tweeviz_api: {{ template "tweeviz.api.fullname" . }}:{{ .Values.port }}
```

There should be an option to disable subchart deployment with
`<subchart>.deployChart` flag and `addresses` map.

In our example Tweeviz UI depends on Tweeviz API and `tweeviz_ui` has the
following in `requirements.yaml`:

```
dependencies:
- name: tweeviz_api
  version: ^0.x
  repository: https://mirantisworkloads.storage.googleapis.com/
  condition: tweeviz_api.deployChart
```

And the following in Values.yaml:
```
tweeviz_api:
  deployChart: false
  addresses:
    tweeviz_api: ""
```