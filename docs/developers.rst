Charts development guide
========================

Repository structure
--------------------

.. code-block::

    k8s-apps
    |_charts
       |_<service1_name>
          |_templates
             |_*.yaml
             |_NOTES.txt
             |_ _helpers.tpl (optional)
          |_Chart.yaml
          |_README.md
          |_values.yaml
          |_requirements.yaml (optional)
       |_<service2_name>
       …
    |_images
       |_<image1_name> (optional)
          |_ Dockerfile
          |_ .version
       |_<image2_name> (optional)
       …
    |_tests
    |_tools
    |_docs


charts directory
----------------

templates
~~~~~~~~~

* Use  `{{ | trunc 63 | trimSuffix "-" }}` in places where the release name is a
  part of the template.
* All k8s object names should have the format `<service-name>-<release-name>`
* All file names should use dashed notation to separate words
  ``some-deployment.yaml``
* Put each k8s object to a separate file and put all service-related files to
  a separate directory.

  Example:

  .. code-block::

      templates/namenode/service.yaml
      templates/namenode/statefulset.yaml
      templates/datanode/service.yaml
      templates/datanode/deployment.yaml

  If there are multiple objects of the same kind under one service (directory),
  they should be prefixed to differentiate them.

* External/end user access - support LoadBalancer, Ingress and NodePort

  * NodePort with ability to specify port

* Persistent storage for data folders - PVC, HostPath and EmptyDir

  * On Public Clouds - PVC
  * For performance on Baremetal - hostPath
  * For local/dev - emptyDir

* `helm.sh/created` annotation should not be used as it has no effect
* All objects should have the following labels:

  .. code-block:: yaml

      labels:
        heritage: "{{ .Release.Service }}"
        release: "{{ .Release.Name }}"
        chart: "{{.Chart.Name}}-{{.Chart.Version}}"
        app: {{ template "fullname" . }}

NOTES.txt
~~~~~~~~~

Should contain:

* Message about service deployment completion
* Addresses of deployed services or a command to get them
* Authentication information (usernames/passwords/etc.)

_helpers.tpl
~~~~~~~~~~~~

* Use for names (in `name:` and `app:` label):

  .. code-block:: smarty

      {{- define "namenode-fullname" -}}
      {{- printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
      {{- end -}}

* Use for addresses if convenient:

  .. code-block:: smarty

      {{- define "namenode-address" -}}
      {{- printf "some address"  | trunc 63 | trimSuffix "-" -}}
      {{- end -}}

* Use for anything that will be referenced more than once

Chart.yaml
~~~~~~~~~~

You should always include the following fields:

.. code-block:: yaml

    description: A Helm chart for Kubernetes
    name: test       # lower case letters, numbers and -
    version: 0.1.0   # initial version 0.1.0

README.md
~~~~~~~~~

* Should have information about how to deploy the chart

values.yaml
~~~~~~~~~~~

* Variables names should begin with a lowercase letter, and words should be
  separated with camelcase (`someParam`: `asd`).
* Try to describe each variable as clearly as possible. Use inline comments
  for that.
* Common variables (should be present in all charts if applicable)

  The following variables should be set per-service in case of multiple
  services. In that case they should be located under `<service-name>` key.

  Image-related variables:

  .. code-block:: yaml

      image:
        repository: mirantisworkloads/zookeeper
        tag: 3.5.3-rc1
        pullPolicy: Always

  Upstream images or images published to `mirantisworkloads` Docker Hub
  repository should be used as defaults.

  Replicas number:

  .. code-block:: yaml

      replicas: 1

  Resources requests and limits. Requests should be set by default, limits are
  optional:

  .. code-block:: yaml

      resources:
        requests:
          cpu: 100m
          memory: 128Mi
        #limits:
          #cpu: 100m
          #memory: 128Mi

  Use `toYaml` function to set them in templates:

  .. code-block:: yaml

          resources:
      {{ toYaml .Values.resources | indent 12 }}


  Persistence configuration:

  .. code-block:: yaml

      persistence:
        type: emptyDir # or hostPath or PersistentVolumeClaim

        #storageClass: ""
        volumeSize: 10Gi

        hostPath: ""

  Java-related variables:

  .. code-block:: yaml

      heapSize: 1G

  Logging-related variables:

  .. code-block:: yaml

      logLevel: INFO

  AntiAffinity-related variables:

  .. code-block:: yaml

      antiAffinity: hard # or soft or null

  Three options should be supported:

  * `hard` - pods will not be scheduled on the same node under any
    circumstances
  * `soft` - pods will not be scheduled on the same node if possible
  * `null` or anything else - antiAffinity will not be used at all

  Should be used in templates in the following way:

  .. code-block:: yaml

      annotations:
        {{ - if eq .Values.antiAffinity "hard" }}
        scheduler.alpha.kubernetes.io/affinity: >
            {
              "podAntiAffinity": {
                "requiredDuringSchedulingIgnoredDuringExecution": [{
                  "labelSelector": {
                    "matchExpressions": [{
                      "key": "app",
                      "operator": "In",
                      "values": ["{{ template "fullname" . }}"]
                    }]
                  },
                  "topologyKey": "kubernetes.io/hostname"
                }]
              }
            }
        {{ - else if eq .Values.antiAffinity "soft" }}
        scheduler.alpha.kubernetes.io/affinity: >
            {
              "podAntiAffinity": {
                "preferredDuringSchedulingIgnoredDuringExecution": [{
                  "weight": 100,
                  "preference": {
                    "matchExpressions": [{
                      "key": "app",
                      "operator": "In",
                      "values": ["{{ template "fullname" . }}"]
                    }]
                  },
                  "topologyKey": "kubernetes.io/hostname"
                }]
              }
            }
        {{ - end }}


  Probe-related variables:

  .. code-block:: yaml

      probeInitialDelaySeconds: 15
      probeTimeoutSeconds: 5

  External/end user access configuration:

  .. code-block:: yaml

      service:
        type: ClusterIP # or NodePort or LoadBalancer

        nodePort: ""

        loadBalancerIP: ""
        loadBalancerSourceRanges: []

        annotations: {}

      ingress:
        enabled: false
        hosts: []
          #   - some.domain
        tls:
          enabled: false
          secretName: ""
        annotations: {}
          #   kubernetes.io/ingress.class: nginx

  External services support:

  .. code-block:: yaml

      some-chart:
        # if disabled, subchart will not be deployed
        deployChart: true
        # this address will be used if subchart deployment is disabled
        externalAddress: ""

  .. NOTE:: it's recommended to move internal/external address selection logic
            to the ``_helpers.tpl``

  Monitoring:

  .. code-block:: yaml

      prometheus:
        jmxExporter:
          enabled: true

requirements.yaml
~~~~~~~~~~~~~~~~~

.. code-block:: yaml

      dependencies:
        - name: some-chart  # name of the chart
          version: ^1.x  # this means >= 1.0.0, < 2
          repository: OUR_REPO_LINK
          condition: some-chart.deployChart

* You can define subchart repository the following way for development
  purposes: `file://../zookeeper`. After making some changes in subcharts,
  you'll have to run `helm dep up` from your chart directory. That way you will
  not have to push charts to repository or create local repository
* You should put `^MAJOR.x` to version field where `MAJOR` is a major version
  of dependent chart and `x` is literally x.
* By default dependant charts will be deployed.
* To use external service you should disable corresponding `deployChart` flag
  and set `externalAddress` instead.

images directory
----------------

* It’s preferred to use upstream images, but if it’s not possible for some
  reasons, images for charts should go there.
* There are no any frameworks to build images. `docker build` is our
  everything.
* Use ARG for templating. https://docs.docker.com/engine/reference/builder/#arg
* Each directory should contain README.md file that describes how to build
  images and which ARGs are supported (if any).
* Each directory should contain ``.version`` file that contains image tag.