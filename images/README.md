
# Chart Images

All images under the charts should be pushed to local repository to prevent bugs and security issues from upstream repos.

In order to automate the process for wide range of existing images the required build/push is managed configured from `DockerMake.yml`.
DockerMake provide predictive, layered docker image builds.

For advanced automation integrate with:

* PyInvoke (https://github.com/pyinvoke/invoke)
* example complex invoke task file: https://github.com/epcim/docker-salt-formulas/blob/master/tasks.py


## Install

We use virtualenv to isolate python requirements to build images.

Prerequsites:

        python -m ensurepip --default-pip
        pip install --user virtualenv
        cd images/
        PYENV=~/.pyenv_dockermake
        python -m virtualenv -p python ${PYENV}

        . ${PYENV}/bin/activate

DockerMake (https://github.com/avirshup/DockerMake):

        pip install DockerMake
        # or for latest
        pip install -e "git+https://github.com/avirshup/DockerMake#egg=dockermake"


## Usage

To list and build targets:

        docker-make --list
        docker-make [target]


Full spec.:

        docker-make [-h] [-f MAKEFILE] [-a] [-l] [--build-arg BUILD_ARG]
                           [--requires [REQUIRES [REQUIRES ...]]] [--name NAME] [-p]
                           [-n] [--dockerfile-dir DOCKERFILE_DIR] [--pull]
                           [--cache-repo CACHE_REPO] [--cache-tag CACHE_TAG]
                           [--no-cache] [--bust-cache BUST_CACHE] [--clear-copy-cache]
                           [--keep-build-tags] [--repository REPOSITORY] [--tag TAG]
                           [--push-to-registry] [--registry-user REGISTRY_USER]
                           [--registry-token REGISTRY_TOKEN] [--version] [--help-yaml]
                           [--debug]
                           [TARGETS [TARGETS ...]]


Workflow:

        TAG="--tag $(date '+%Y%m%d')"
        PUSH="--push-to-registry -u docker.io/mirantisworkloads --registry-user push --registry-token push"
        CACHE="--no-cache"

        docker-make --list
        docker-make [TARGET] $TAG
        docker-make --all    $TAG $PUSH_OPTS

        # to build all targets (that has FROM* specified):
        docker-make --list |awk '/*/{print $2}'| grep -v 'base\|common' | xargs -n1 -I% docker-make % $TAG $PUSH $CACHE


NOTE: https://github.com/avirshup/DockerMake/issues/52 (to auto-ignore targets without FROM spec.)




