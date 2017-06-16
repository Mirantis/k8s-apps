Verify applications work on new k8s or helm version
===================================================

Every time when new kubernetes or helm version become available, important to
verify that all charts are work correctly with different install parameters.

Requirements
------------

In order to fully verify new version of k8s/helm, host should have:

* at least one available StorageClass (for example, glusterfs);
* kubernetes ingress controller.

Important items to verify
-------------------------

The following steps should be taken in order to verify that applications are working correctly:

#. Install at least one chart, supporting PV, with enabled *persistence* values,
   for example:

   .. code-block:: yaml

      persistence:
        type: PersistentVolumeClaim
        storageClass: glusterfs
        volumeSize: 10Gi

#. Install at least one chart, supporting antiAffinity, with enabled
   *antiAffinity* value ("hard" and "soft").

#. Install at least one chart with different service types, for example:

   .. code-block:: yaml

      service:
        type: "NodePort" # or LoadBalancer

#. Install at least one chart with dependencies, for example `kafka`, which
   depends on `zookeeper`.

#. Ingress controller should be installed. Verify ingress resources with deploying
   at least one chart with *ingress* values:

   .. code-block:: yaml

      ingress:
        enabled: true
        hosts: [<some-host>]
        tls:
          enabled: true
          secretName: "<some-secretName>"

#. Important to check scalability: choose at least one chart with scaling
   support and scale up it with command:

   :code:`kubectl scale <resource> --replicas=<replicas number>`

#. Delete chart and verify that all it's resources are successfully deleted.

#. All charts should be installed without any errors with default values. It
   could be done with k8s-apps tests:

   :code:`go run tools/helm-test/main.go --exclude tweepub,tweetics`

#. Run one of the available scenarios:

   ::

      cd scenarios/tweeanalytics/
      bash tweeanalytics.sh --app-key <TWITTER_APP_KEY>\
                            --app-secret <TWITTER_APP_SECRET>\
                            --token-key <TWITTER_TOKEN_KEY>\
                            --token-secret <TWITTER_TOKEN_SECRET>\
                            [-cas][--mode]

#. Check standard helm commands:

   * :code:`helm install ... --dry-run --debug` will show rendered templates.
   * :code:`helm lint <chart>` will run linter for specified chart.
   * :code:`helm repo add mirantisworkloads https://mirantisworkloads.storage.googleapis.com`
     will add mirantisworkloads repo to current helm.
   * :code:`helm list` will list current installed charts releases.
   * :code:`helm upgrade <release> <chart>` to in-place upgrade template without re-installing chart.
   * :code:`helm dep up` update chart's dependencies in accordance with requirements.yaml file.
