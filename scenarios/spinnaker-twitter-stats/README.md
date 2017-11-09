# Spinnaker pipelines for Twitter Stats

## Installation

* Install tiller in kube-system namespace:
  ```shell
  helm init --tiller-namespace kube-system
  ```
* Add mirantisworkloads repository:
  ```shell
  helm repo add mirantisworkloads https://mirantisworkloads.storage.googleapis.com
  ```
* Set twitter credentials in `spinnaker-values.yaml`:

  Replace CHANGEME with proper values.

  ```shell
    export TS_APP_KEY=CHANGEME
    export TS_APP_SECRET=CHANGEME
    export TS_TOKEN_KEY=CHANGEME
    export TS_TOKEN_SECRET=CHANGEME
  ```

* Install spinnaker chart with pre-uploaded pipelines:
  ```shell
  helm install mirantisworkloads/spinnaker -f spinnaker-values.yaml
  ```

## Usage

* Go to the spinnaker UI and move to `twitter-analytics` application.
* Run `Twitter Analytics: deploy` pipeline.

  It has only one parameter - `Name` - that will be used as namespace and
  prefix for helm releases.

  Deploy pipeline will start twitter-stats services deployment. One service
  per stage will be deployed.

* Run `Twitter Analytics: destroy` pipeline with the same `Name` for cleanup.

Please refer to the [twitter-stats](../twitter-stats/README.md) for more information.