local k = import "ksonnet.beta.2/k.libsonnet";

local kStateful = k.apps.v1beta1.statefulSet;
local kAntiAffinity = kStateful.mixin.spec.template.spec.affinity.podAntiAffinityType;
local kContainer = k.apps.v1beta1.statefulSet.mixin.spec.template.spec.containersType;
local kProbe = kContainer.mixin.livenessProbeType;
local kVolume =  kStateful.mixin.spec.template.spec.volumesType;

{
  platform:: {
    statefulSet:: {
      new(serviceName, replicas)::
        kStateful.new() +
        kStateful.mixin.spec.replicas(replicas) +
        kStateful.mixin.spec.serviceName(serviceName),
      setMeta(name, labels={}, annotations={})::
        kStateful.mixin.metadata.name(name) +
        (if std.length(labels) > 0 then kStateful.mixin.metadata.labels(labels) else {}) +
        kStateful.mixin.spec.template.metadata.name(name) +
        (if std.length(labels) > 0 then kStateful.mixin.spec.template.metadata.labels(labels) else {}) +
        (if std.length(annotations) > 0 then kStateful.mixin.spec.template.metadata.annotations(annotations) else {}),
      antiAffinity(type, name)::
        (if type != "null" then kStateful.mixin.spec.template.spec.affinity.podAntiAffinity.mixinInstance($.platformUtils.antiAffinity(type, name)) else {}),
      initContainers(initContainers, expr=true)::
        if expr then kStateful.mixin.spec.template.spec.initContainers(initContainers) else {},
      containers(containers)::
        if std.length(containers) > 0 then
          kStateful.mixin.spec.template.spec.containers(containers)
        else {},
      volumes(type, name, hostPath="", storageClass="", volumeSize="10Gi")::
        if type == "emptyDir" then
          kStateful.mixin.spec.template.spec.volumes(kVolume.fromEmptyDir(name))
        else if type == "hostPath" then
          kStateful.mixin.spec.template.spec.volumes(kVolume.fromHostPath(name, hostPath))
        else if type == "PersistentVolumeClaim" then
         kStateful.mixin.spec.volumeClaimTemplates(
           kStateful.mixin.spec.volumeClaimTemplatesType.mixin.metadata.name(name) +
           (if storageClass != "" then kStateful.mixin.spec.volumeClaimTemplatesType.mixin.metadata.annotations({"volume.beta.kubernetes.io/storage-class": storageClass}) else {}) +
           kStateful.mixin.spec.volumeClaimTemplatesType.mixin.spec.accessModes("ReadWriteOnce") +
           {spec+: {resources: {requests: {storage: volumeSize}}}}),
      mixinVolumes(volumes)::
        if std.length(volumes) > 0 then
          kStateful.mixin.spec.template.spec.volumes(volumes)
        else {}
    },
    container:: {
      new(name, image, pullPolicy="IfNotPresent"):: kContainer.new(name, image) + kContainer.imagePullPolicy(pullPolicy),
      readinessProbe(readinessProbe, initDelay, timeout)::
          readinessProbe +
          kContainer.mixin.readinessProbe.initialDelaySeconds(initDelay) +
          kContainer.mixin.readinessProbe.timeoutSeconds(timeout),
      livenessProbe(livenessProbe, initDelay, timeout)::
          livenessProbe +
          kContainer.mixin.livenessProbe.initialDelaySeconds(initDelay) +
          kContainer.mixin.livenessProbe.timeoutSeconds(timeout),
      logVolumeMount(path, enabled=false):: $.agents.logging.volumeMount.new(path, enabled)
    }
  },
  agents:: {
    logging:: {
      volumeMount:: {
        new(path, enabled=false)::
          if enabled then
            kContainer.volumeMounts(kContainer.volumeMountsType.new("logdir", path))
          else {}
      },
      volume:: {
        new(name, enabled=false)::
          if enabled then
            [kVolume.fromConfigMap("filebeat-config", name, [{key: "filebeat-yml", path: "filebeat.yml"}]),
             kVolume.fromEmptyDir("logdir")]
          else {}
      },
      container:: {
        new(logdirPath, enabled=false)::
          if enabled then
            kContainer.new("filebeat", "mirantisworkloads/filebeat:1.0.0") +
            kContainer.imagePullPolicy("IfNotPresent") +
            kContainer.volumeMounts(kContainer.volumeMountsType.new("filebeat-config", "/etc/filebeat")) +
            $.agents.logging.volumeMount.new(logdirPath, enabled) +
            kContainer.command(["filebeat", "-c", "/etc/filebeat/filebeat.yml", "-e", "-d", "\"*\""])
          else {}
      }
    }
  },
  platformUtils:: {
    antiAffinity(type="soft", fullname)::
      if type == "hard" then
        kAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution([{
          "labelSelector": {
            "matchExpressions": [{
              "key": "app",
              "operator": "In",
              "values": [fullname]
            }]
          },
          "topologyKey": "kubernetes.io/hostname"
        }])
      else if type == "soft" then
        kAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution([{
          "weight": 100,
          "preference": {
            "matchExpressions": [{
              "key": "app",
              "operator": "In",
              "values": [fullname]
            }]
          },
          "topologyKey": "kubernetes.io/hostname"
        }])
      else {},
  }
}
