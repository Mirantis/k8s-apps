**Using**

1. Use contest "service-catalog" to communicate with the service catalog API server:

```bash
$ kubectl config use-context service-catalog
```

2. Build using make:

```bash
$ make build
```

3. Run script:

```bash
$ ./service-catalog-client
```
You can override the following params:

```bash
$ ./service-catalog-client --broker broker --instance instance --chart zookeeper --version 1.1.0 --namespace test-ns --binding binding
```
Also you can override chart values in the values.json file.