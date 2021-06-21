## Description

This plugin enables you to manage artifacts in [Artifactory](https://jfrog.com/artifactory/) in a Vela pipeline.

Source Code: https://github.com/go-vela/vela-artifactory

Registry: https://hub.docker.com/r/target/vela-artifactory

## Usage

> **NOTE:**
>
> Users should refrain from using latest as the tag for the Docker image.
>
> It is recommended to use a semantically versioned tag instead.

Sample of copying an artifact:

```yaml
steps:
  - name: copy_artifacts
    image: target/vela-artifactory:latest
    pull: always
    parameters:
      action: copy
      path: libs-snapshot-local/foo.txt
      target: libs-snapshot-local/bar.txt
      url: http://localhost:8081/artifactory
```

Sample of deleting an artifact:

```yaml
steps:
  - name: delete_artifacts
    image: target/vela-artifactory:latest
    pull: always
    parameters:
      action: delete
      path: libs-snapshot-local/foo.txt
      url: http://localhost:8081/artifactory
```

Sample of setting properties on an artifact:

```yaml
steps:
  - name: set_properties_artifacts
    image: target/vela-artifactory:latest
    pull: always
    parameters:
      action: set-prop
      path: libs-snapshot-local/foo.txt
      props:
        - name: single
          value: foo
        - name: multiple
          values:
            - bar
            - baz
      url: http://localhost:8081/artifactory
```

Sample of uploading an artifact:

```yaml
steps:
  - name: upload_artifacts
    image: target/vela-artifactory:latest
    pull: always
    parameters:
      action: upload
      path: libs-snapshot-local/
      sources:
        - foo.txt
        - target/*.jar
        - dist/**/*.js
      url: http://localhost:8081/artifactory
```

Sample of pretending to upload an artifact:

```diff
steps:
  - name: upload_artifacts
    image: target/vela-artifactory:latest
    pull: always
    parameters:
      action: upload
+     dry_run: true
      path: libs-snapshot-local/
      sources:
        - foo.txt
        - target/*.jar
        - dist/**/*.js
      url: http://localhost:8081/artifactory
```

Sample of using docker-promote on an artifact:

```yaml
steps:
  - name: docker_promote_artifacts
    image: target/vela-artifactory:latest
    pull: always
    parameters:
      action: docker-promote
      target_repo: libs-snapshot-local
      docker_registry: octocat/hello-world
      tag: latest
      target_docker_registry: octocat/hello-world
      target_tags: "${VELA_BUILD_COMMIT:0:8}"
```

## Secrets

> **NOTE:** Users should refrain from configuring sensitive information in your pipeline in plain text.

### Internal

The plugin accepts the following `parameters` for authentication:

| Parameter   | Environment Variable Configuration           |
| ----------- | -------------------------------------------- |
| `api_key`   | `PARAMETER_API_KEY`, `ARTIFACTORY_API_KEY`   |
| `password`  | `PARAMETER_PASSWORD`, `ARTIFACTORY_PASSWORD` |
| `username`  | `PARAMETER_USERNAME`, `ARTIFACTORY_USERNAME` |

Users can use [Vela internal secrets](https://go-vela.github.io/docs/tour/secrets/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: copy_artifacts
    image: target/vela-artifactory:latest
    pull: always
+   secrets: [ artifactory_username, artifactory_password ]
    parameters:
      action: copy
      path: libs-snapshot-local/foo.txt
      target: libs-snapshot-local/bar.txt
      url: http://localhost:8081/artifactory
-     username: octocat
-     password: superSecretPassword
```

> This example will add the secrets to the `copy_artifacts` step as environment variables:
>
> * `ARTIFACTORY_USERNAME=<value>`
> * `ARTIFACTORY_PASSWORD=<value>`

### External

The plugin accepts the following files for authentication:

| Parameter  | File Configuration                                                            |
| ---------- | ----------------------------------------------------------------------------- |
| `api_key`  | `/vela/parameters/artifactory/api_key`, `/vela/secrets/artifactory/api_key`   |
| `password` | `/vela/parameters/artifactory/password`, `/vela/secrets/artifactory/password` |
| `username` | `/vela/parameters/artifactory/username`, `/vela/secrets/artifactory/username` |

Users can use [Vela external secrets](https://go-vela.github.io/docs/tour/secrets/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: copy_artifacts
    image: target/vela-artifactory:latest
    pull: always
    parameters:
      action: copy
      path: libs-snapshot-local/foo.txt
      target: libs-snapshot-local/bar.txt
      url: http://localhost:8081/artifactory
-     username: octocat
-     password: superSecretPassword
```

> This example will read the secret values in the build workspace stored at `/vela/secrets/artifactory/*`

## Parameters

> **NOTE:**
>
> The plugin supports reading all parameters via:
>
> * environment variables - `PARAMETER_*` OR `ARTIFACTORY_*`
> * files - `/vela/parameters/artifactory/*` OR `/vela/secrets/artifactory/*`
>
> Any values set from a file takes precedence over values set from the environment.

The following parameters are used to configure the image:

| Name        | Description                                  | Required | Default | Environment Variables                            |
| ----------- | -------------------------------------------- | -------- | ------- | ------------------------------------------------ |
| `action`    | action to perform against Artifactory        | `true`   | `N/A`   | `PARAMETER_ACTION`<br>`ARTIFACTORY_ACTION`       |
| `api_key`   | API key for communication with Artifactory   | `false`  | `N/A`   | `PARAMETER_DRY_RUN`<br>`ARTIFACTORY_DRY_RUN`     |
| `dry_run`   | enables pretending to perform the action     | `false`  | `false` | `PARAMETER_API_KEY`<br>`ARTIFACTORY_API_KEY`     |
| `log_level` | set the log level for the plugin             | `true`   | `info`  | `PARAMETER_LOG_LEVEL`<br>`ARTIFACTORY_LOG_LEVEL` |
| `password`  | password for communication with Artifactory  | `false`  | `N/A`   | `PARAMETER_PASSWORD`<br>`ARTIFACTORY_PASSWORD`   |
| `url`       | Artifactory instance to communicate with     | `true`   | `N/A`   | `PARAMETER_URL`<br>`ARTIFACTORY_URL`             |
| `username`  | user name for communication with Artifactory | `false`  | `N/A`   | `PARAMETER_USERNAME`<br>`ARTIFACTORY_USERNAME`   |

### Copy

The following parameters are used to configure the `copy` action:

| Name        | Description                                         | Required | Default | Environment Variables                            |
| ----------- | --------------------------------------------------- | -------- | ------- | ------------------------------------------------ |
| `flat`      | enables removing source directory hierarchy         | `false`  | `false` | `PARAMETER_FLAT`<br>`ARTIFACTORY_FLAT`           |
| `path`      | source path to copy artifact(s) from                | `true`   | `N/A`   | `PARAMETER_PATH`<br>`ARTIFACTORY_PATH`           |
| `recursive` | enables copying sub-directories for the artifact(s) | `false`  | `false` | `PARAMETER_RECURSIVE`<br>`ARTIFACTORY_RECURSIVE` |
| `target`    | target path to copy artifact(s) to                  | `true`   | `N/A`   | `PARAMETER_TARGET`<br>`ARTIFACTORY_TARGET`       |

### Delete

The following parameters are used to configure the `delete` action:

| Name        | Description                                          | Required | Default | Environment Variables                            |
| ----------- | ---------------------------------------------------- | -------- | ------- | ------------------------------------------------ |
| `path`      | target path to delete artifact(s) from               | `true`   | `N/A`   | `PARAMETER_PATH`<br>`ARTIFACTORY_PATH`           |
| `recursive` | enables removing sub-directories for the artifact(s) | `false`  | `false` | `PARAMETER_RECURSIVE`<br>`ARTIFACTORY_RECURSIVE` |

### Docker-Promote

The following parameters are used to configure the `docker-promote` action:

| Name                     | Description                                         | Required | Default | Environment Variables                                                      |
| ------------------------ | --------------------------------------------------- | -------- | ------- | -------------------------------------------------------------------------- |
| `copy`                   | set to copy instead of moving the image             | `false`  | `true`  | `PARAMETER_COPY`<br>`ARTIFACTORY_COPY`                                     |
| `docker_registry`        | path to image in docker registry                    | `true`   | `N/A`   | `PARAMETER_DOCKER_REGISTRY`<br>`ARTIFACTORY_DOCKER_REGISTRY`               |
| `promote_props`          | enables setting properties on the promoted artifact | `false`  | `false` | `PARAMETER_PROMOTE_PROPS`<br>`ARTIFACTORY_PROMOTE_PROPS`                   |
| `tag`                    | name of the tag for promoting                       | `true`   | `N/A`   | `PARAMETER_TAG`<br>`ARTIFACTORY_TAG`                                       |
| `target_docker_registry` | path for target image in docker registry            | `true`   | `N/A`   | `PARAMETER_TARGET_DOCKER_REGISTRY`<br>`ARTIFACTORY_TARGET_DOCKER_REGISTRY` |
| `target_repo`            | name of the docker registry containing the image    | `true`   | `N/A`   | `PARAMETER_TARGET_REPO`<br>`ARTIFACTORY_TARGET_REPO`                       |
| `target_tags`            | name of the final tags after promotion              | `true`   | `N/A`   | `PARAMETER_TARGET_TAGS`<br>`ARTIFACTORY_TARGET_TAGS`                       |

### Set-Prop

The following parameters are used to configure the `set-prop` action:

| Name    | Description                          | Required | Default | Environment Variables                    |
| ------- | ------------------------------------ | -------- | ------- | ---------------------------------------- |
| `path`  | target path to artifact(s)           | `true`   | `N/A`   | `PARAMETER_PATH`<br>`ARTIFACTORY_PATH`   |
| `props` | properties to set on the artifact(s) | `true`   | `N/A`   | `PARAMETER_PROPS`<br>`ARTIFACTORY_PROPS` |

### Upload

The following parameters are used to configure the `upload` action:

| Name           | Description                                           | Required | Default | Environment Variables                                  |
| -------------- | ----------------------------------------------------- | -------- | ------- | ------------------------------------------------------ |
| `flat`         | enables removing source directory hierarchy           | `false`  | `false` | `PARAMETER_FLAT`<br>`ARTIFACTORY_FLAT`                 |
| `include_dirs` | enables including sub-directories for the artifact(s) | `false`  | `false` | `PARAMETER_INCLUDE_DIRS`<br>`ARTIFACTORY_INCLUDE_DIRS` |
| `path`         | target path to upload artifact(s) to                  | `true`   | `N/A`   | `PARAMETER_PATH`<br>`ARTIFACTORY_PATH`                 |
| `recursive`    | enables uploading sub-directories for the artifact(s) | `false`  | `false` | `PARAMETER_RECURSIVE`<br>`ARTIFACTORY_RECURSIVE`       |
| `regexp`       | enables reading the sources as a regular expression   | `false`  | `false` | `PARAMETER_REGEXP`<br>`ARTIFACTORY_REGEXP`             |
| `sources`      | list of artifact(s) to upload                         | `true`   | `N/A`   | `PARAMETER_SOURCES`<br>`ARTIFACTORY_SOURCES`           |

## Template

COMING SOON!

## Troubleshooting

You can start troubleshooting this plugin by tuning the level of logs being displayed:

```diff
steps:
  - name: copy_artifacts
    image: target/vela-artifactory:latest
    pull: always
    parameters:
      action: copy
+     log_level: trace
      path: libs-snapshot-local/foo.txt
      target: libs-snapshot-local/bar.txt
      url: http://localhost:8081/artifactory
```

Below are a list of common problems and how to solve them:
