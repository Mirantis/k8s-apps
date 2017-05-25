# ELK cluster

## Deploy chart
```console
$ helm install .
```

To run ELK helm tests:

    helm test {{ .Release.Name }} --cleanup
