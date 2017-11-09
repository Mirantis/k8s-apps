# Jenkins job for Twitter Stats

## Installation

* Install tiller in kube-system namespace:
  ```shell
  helm init --tiller-namespace kube-system
  ```
* Add mirantisworkloads repository:
  ```shell
  helm repo add mirantisworkloads https://mirantisworkloads.storage.googleapis.com
  ```
* Set twitter credentials in `jenkins-values.yaml`:

  Replace CHANGEME with proper values.

    export TS_APP_KEY=CHANGEME
    export TS_APP_SECRET=CHANGEME
    export TS_TOKEN_KEY=CHANGEME
    export TS_TOKEN_SECRET=CHANGEME

* Install jenkins chart with pre-defined job:
  ```shell
  helm install mirantisworkloads/jenkins -f jenkins-values.yaml
  ```

## Usage

Go to the jenkins UI and run `twitter` job.

It has the following parameters:
 * TS_NAME   - namespace and prefixes for helm releases
 * TS_CMD    - `up`, `test` or `down` (can be combined)
 * TS_CHARTS - list of charts to process (`zookeeper hdfs kafka spark tweepub tweetics tweeviz` if not set)

Jenkins executor will be created as a pod in the same Kubernetes cluster.

Please refer to the [twitter-stats](../twitter-stats/README.md) for more information.