# Fission chart

Fission is a fast serverless framework for Kubernetes with a focus on developer
productivity and high performance.

Fission is a FaaS (Function as a Service) - users create functions (source level),
register them with fission using a CLI, and associate functions with triggers.

## Installation

### Fission service

TODO: implement helm chart to deploy Fission service.

### Fission client

To install fission client on host, just run:

```
curl -Lo fission https://github.com/fission/fission/releases/download/nightly20170621/fission-cli-linux && chmod +x fission && sudo mv fission /usr/local/bin/
```

### Fission UI

TODO: Fission has community maintained UI. To run fission UI, install chart with
enabled UI value.

## Configuration

TODO: Add values table with description and defaults.

## General concepts

Fission operates with **functions** - a peace of code that follows the fission
interface.

Functions runs on **environment containers** - language- and runtime-specific
containers. Fission comes with NodeJS and Python environments and could be
extended with custom environments (an environment is essentially just a
container with a webserver and dynamic loader).

A **trigger** is something that maps an event to a function. Only *HTTP routes*
supported today.

## Architecture

Fission consists of 4 neutral-language components:

* **controller** - stateful component; contains CRUD APIs for functions, http
  triggers, environments, Kubernetes event watches.  It needs to be configured
  with a URL to an etcd cluster and a path to a persistent volume. The volume is
  used to store the functions' source code. Etcd is used as the DB.

* **poolmgr** manages pools of generic containers and function containers;
  watches the controller API and eagerly creates generic pools for environments.

* **router** - forwards HTTP requests to function pods. If there's no running
  service for a function, it requests one from poolmgr, while holding on to
  the request.

* **kubewatcher** - watches the Kubernetes API and invokes functions associated
  with watches, sending the watch event to the function.

Also it contains language-specific containers, called
**environment containers**. These contianers runs user-defined functions. They
must contain an HTTP server and a loader for functions.

Poolmgr deploys the environment container into a pod with fetcher. When poolmgr
needs to create a service for a function, it calls fetcher to fetch the
function. Fetcher downloads the function into a volume shared between fetcher
and this environment container. Poolmgr then requests the container to load the
function.

Fission has logger, that helps to forward function logs to centralized db
service for log persistence. Currently only *influxdb* is supported to store
logs.

## Simple usage

```
# Add the stock NodeJS env to your Fission deployment
$ fission env create --name nodejs --image fission/node-env

# A javascript one-liner that prints "hello world"
$ curl https://raw.githubusercontent.com/fission/fission/master/examples/nodejs/hello.js > hello.js

# Upload your function code to fission
$ fission function create --name hello --env nodejs --code hello.js

# Map GET /hello to your new function
$ fission route create --method GET --url /hello --function hello

# Run the function.  This takes about 100msec the first time.
$ curl http://$FISSION_ROUTER/hello
Hello, world!
```

## Supported environments

Following environments come with Fission:

* Python (*fission/python-env* and *fission/python3-env*)

* NodeJS (*fission/node-env*)

* Golang (*fission/go-env*)

More environments implemented here:
https://github.com/fission/fission/tree/master/environments, but they aren't
pushed to docker hub.

## Future features

Fission roadmap has many awesome useful features, which doesn't implemented yet.
To learn it, please visit: https://github.com/fission/fission/blob/master/Documentation/Roadmap.md.
