# Rollout Chart


## Setup

The demo is designed to be executed on two Kubernetes clusters (and
even three if chart components are deployed separately), however it can
be launched on a single cluster if desired.

1. Add mirantisworkloads helm repository:

    ```console
    $ helm repo add mirantisworkloads https://mirantisworkloads.storage.googleapis.com
    ```

2. Prepare configuration file (replace all CHANGEME occurrences in values.yaml):

    * Set proper credentials for your twitter account:

        ```yaml
        appKey: CHANGEME
        appSecret: CHANGEME
        tokenKey: CHANGEME
        tokenSecret: CHANGEME
        ```

    * Configure docker repository and credentials:

      ```yaml
      docker:
        user: CHANGEME
        password: CHANGEME
        repository: CHANGEME
        stableTag: CHANGEME
      ```

      Set repository for tweeviz-ui:
      ```yaml
      releases:
          tweeviz-ui-${ENV}:
            chart: mirantis/tweeviz_ui
            cluster: frontend
            wait: true
            dependencies:
              tweeviz_api: tweeviz-api-${ENV}
            parameters:
              image:
                repository: CHANGEME/
                tag: v${VERSION}
      ```

    * It's mandatory to have `frontend` and `backend` files under `kubeConfigs` parameter.
      And it's also mandatory to have `frontend` and `backend` contexts defined inside of them.

    * You should define proper `externalIP` for both `clusters` under `helmApplyConfig`
      group.

    * If your cluster has RBAC enabled, the following should be added:
      ```yaml
      rbac:
        enabled: true

      spinnaker:
        rbac:
          enabled: true
        jenkins:
          Master:
            rbac:
              enabled: true
      ```
    * You can optionally expose whatever services you want. They will have
      `NodePort` type by default. For more details about spinnaker and
      gerrit configuration please refer to their charts.

3. Install rollout chart (supply it with previously prepared config file):

    ```console
    $ helm install mirantisworkloads/rollout -f values.yaml
    ```
