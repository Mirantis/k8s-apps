====================================
Charts' log collector implementation
====================================

Current implementation of **logstash** allows to collect logs from any
application in the repo with adding logstash agent container to it's chart.
Logstash agent named **filebeat**. To learn more, how **filebeat** works,
please refer to
`corresponding document page <https://www.elastic.co/guide/en/beats/filebeat/5.2/how-filebeat-works.html>`_.

To activate log collector support in a chart, please read step by step guide
below.

Add filebeat to chart
~~~~~~~~~~~~~~~~~~~~~

#. Add next dict to ``values.yaml`` file:

   .. code-block:: yaml

      logCollector:
        enabled: false
        image:
          repository: mirantisworkloads/
          name: filebeat
          tag: 5.2.2
          pullPolicy: Always
        logstashHosts: [] # external logstash hosts with ports; override with actual data

#. Add `filebeat` configmap to :file:`templates/` dir similar with this
   :file:`charts/kafka/templates/filebeat-configmap.yaml`.

#. Append new container with filebeat application to template(s), where needed
   to collect logs:

   .. code-block:: yaml

      {{- if .Values.logCollector.enabled }}
      - name: filebeat
        image: "{{ .Values.logCollector.image.repository }}{{ .Values.logCollector.image.name }}:{{ .Values.logCollector.image.tag }}"
        imagePullPolicy: {{ .Values.logCollector.image.pullPolicy | quote }}
        volumeMounts:
        - name: filebeat-config
          mountPath: /etc/filebeat
        # and then add logs dirs, shared between container with logs and filebeat container
        - name: logdir
          mountPath: /some/path
        ...
        command:
          - "filebeat"
          - "-c"
          - "/etc/filebeat/filebeat.yml"
          - "-e"
          - "-d"
          - "\"*\""
      {{- end }}

#. Add ``volumes`` entries for shared logs directories between the container
   and filebeat:

   .. code-block:: yaml

      {{- if .Values.logCollector.enabled }}
      - name: filebeat-config
        configMap:
          name: {{ printf "kafka-fb-%s" .Release.Name | trunc 63 }}
          items:
            - key: filebeat-yml
              path: filebeat.yml
      # and again logs dirs, shared between container with logs and filebeat container
      - name: logdir
        emptyDir: {}
      ...
      {{- end }}

#. Add ``volumeMounts`` entries to the actual container spec, which will
   contain dirs with logs, which should be sent to logstash (for example,
   "/opt/kafka/logs").

#. Now all chart changes are finished, time to build image and run application!
   Build and push :file:`image/filebeat/` image and install helm chart.

.. tip:: If there're errors during deploy, try to debug what has been done.

After finishing all these steps, defined logs will be sent to `logstash`
deployment and saved to `elasticsearch` index in json-type without any
filtering and formatting with name *<string-type component name>-MM.dd.YYYY*.

.. tip:: Don't forget to deploy `logstash` and `elasticsearch` to enable
         log collector!

To add filter and format, it's necessary to add new entry to `filter` section
of :file:`chart/logstash/templates/logstash-configmap.yaml` configmap.

Add filter entry to `logstash` configmap
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Add next condition to logstash config filter section:

::

  else if ([fields][log_host] == <string-type component name>) {
      ...
  }

Detailed information about current filters could be found on
`this document page <https://www.elastic.co/guide/en/logstash/current/filter-plugins.html>`_.

Very important to add at least next filter plugin, because it's necessary for
correct elasticsearch index name:

::

  grok {
      match => {
          "message" => ["%{GREEDYDATA:message}"]
      }
      overwrite => [ "message" ]
      add_field => {
          "received_from" => "%{host}"
      }
  }

Elasticsearch index name will be built with field :code:`%{received_from}`.

Now check index records and ensure that filter is working.

.. tip:: If not, try to debug what is wrong.

Conclusion
~~~~~~~~~~

Now the chart supports logs collector, which passed to storage. Play with
filters and `logstash` configmap for the best result.
