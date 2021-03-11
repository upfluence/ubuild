# ubuild

ubuild is our release tool. It allows you to:

    Generate a release
    Compile programs
    Build containers
    Push those containers to a remote registry

Everything in one command.

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
custom_bump_strategies:
   <branch>: <strategie> # valid strategies are "bump_patch", "bump_rc"
```

For example this configuration will build an image based on a python package, update the release of upfluence/ner_analyser and update to the docker image named "ner-analyser".
```yaml
type: py
verbose: true
repository: <org>/<repository>
docker:
  image: "<image_name>"
```

## Circle CI
To use ubuild with circle-ci you must use the golang primary image have run go get  in a step of you job. The sensitive environment variable should be fetched from the project settings.

Example:    - run: docker images

```yaml
build:
  name: Build

  runs-on: ubuntu-latest

  needs: test

  steps:
    - name: Checkout
      uses: actions/checkout@v2

    - name: Install ubuild
      run: |
        curl -sSL https://github.com/upfluence/ubuild/releases/download/v0.2.0/ubuild-linux-amd64-0.2.0 > ~/go/bin/ubuild
        chmod +x ~/go/bin/ubuild

    - name: Build
      run: |
        echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
        ubuild
      env:
        DOCKER_USERNAME: ${{ secrets.DOCKER_USERNAME }}
        DOCKER_PASSWORD: ${{ secrets.DOCKER_PASSWORD }}
        RELEASE: true
        GITHUB_TOKEN: ${{ secrets.PAT_TOKEN }}
```
