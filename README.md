# ubuild

Ubuild is our docker image builder utilities. It packs several tools and workflow that we use through most of our projects. Ubuild objective is two fold, build the docker image and push it to docker hub.

## Install

`go get  github.com/upfluence/ubuild`

## Prerequisite

Before running ubuild  you must export a GITHUB_TOKEN and set RELEASE to true  into your environment and be logged in into a docker hub account with access to upfluence image registry.

## Configuration
```yaml
type: <lang> go, rb, py, frontend
verbose: <true/false>
repository: <github_path>
# needed only for compiled language
compiler:
  binaries:
    - path: <path_to_entrypp>
  dist:
  CGO: 
  args: # map of arguments to pass to the compiler
    key: val
    ... 
docker:
  dockerfile: <path_to_dockerfile> # Dockerfile by default
  image: <image_name>
  tags: # additional tags
    key: val
    ...
deployer:
  url:
  envs:
    key: val
    ...
```

For example this configuration will build an image based on a python package, update the release of upfluence/ner_analyser and update to the docker image named "ner-analyser".
```yaml
type: py
verbose: true
repository: <org>/<repositoryW
docker:
  image: "<image_name>"
```

## Circle CI
To use ubuild with circle-ci you must use the golang primary image have run go get  in a step of you job. The sensitive environment variable should be set in the project settings.

Example:    - run: docker images

```yaml
publish:
  docker:
    - image: circleci/golang:1.13

  environment:
    RELEASE: true

  steps:
    - checkout

    - setup_remote_docker
    - run:
        command: |
          docker login -u $DOCKER_USER -p $DOCKER_PASS

    - run: go get github.com/upfluence/ubuild/cmd/ubuild

    - run:
        name: Building server image
        command: ubuild
```
