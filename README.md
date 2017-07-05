# Mirantis Platform Services

## Supported helm and kubernetes versions

Currently used kubernetes version **v1.5** with helm version **v2.3.1**. Stable
work is guaranteed.

Next versions for stable work are:

 * kubernetes version **v1.6**
 * helm version **v2.4.2**

These versions are tested together and approved with successful verify version
check list (https://github.com/Mirantis/k8s-apps/blob/master/docs/check_helm_or_k8s_version.rst).

*NOTE*: Do not use helm version **v2.5.0** - chart dependencies are broken in
it. Fix will be landed in **v2.5.1**.
