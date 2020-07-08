## Description

This plugin enables you to manage artifacts in [Artifactory](https://jfrog.com/artifactory/) in a Vela pipeline.

Source Code: https://github.com/go-vela/vela-artifactory

Registry: https://hub.docker.com/r/target/vela-artifactory

## Usage

Sample of copying an artifact:

```yaml
steps:
  - name: copy_artifacts
    image: target/vela-artifactory:v0.2.0
    pull: true
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
    image: target/vela-artifactory:v0.2.0
    pull: true
    parameters:
      action: delete
      path: libs-snapshot-local/foo.txt
      url: http://localhost:8081/artifactory
```

Sample of setting properties on an artifact:

```yaml
steps:
  - name: set_properties_artifacts
    image: target/vela-artifactory:v0.2.0
    pull: true
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
    image: target/vela-artifactory:v0.2.0
    pull: true
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
    image: target/vela-artifactory:v0.2.0
    pull: true
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
    image: target/vela-artifactory:v0.2.0
    pull: true
    parameters:
      action: docker-promote
      target_repo: libs-snapshot-local
      docker_registry: octocat/hello-world
      tag: latest
      target_docker_registry: octocat/hello-world
      target_tags: "${BUILD_COMMIT:0:8}"
```

## Secrets

**NOTE: Users should refrain from configuring sensitive information in your pipeline in plain text.**

You can use Vela secrets to substitute sensitive values at runtime:

```diff
steps:
  - name: copy_artifacts
    image: target/vela-artifactory:v0.2.0
    pull: true
+   secrets: [ artifactory_username, artifactory_password ]
    parameters:
      action: copy
      path: libs-snapshot-local/foo.txt
      target: libs-snapshot-local/bar.txt
      url: http://localhost:8081/artifactory
-     username: octocat
-     password: superSecretPassword
```

## Parameters

The following parameters are used to configure the image:

| Name        | Description                                  | Required | Default |
| ----------- | -------------------------------------------- | -------- | ------- |
| `action`    | action to perform against Artifactory        | `true`   | `N/A`   |
| `api_key`   | API key for communication with Artifactory   | `false`  | `N/A`   |
| `dry_run`   | enables pretending to perform the action     | `false`  | `false` |
| `log_level` | set the log level for the plugin             | `true`   | `info`  |
| `password`  | password for communication with Artifactory  | `false`  | `N/A`   |
| `path`      | source/target path to artifact(s) for action | `true`   | `N/A`   |
| `url`       | Artifactory instance to communicate with     | `true`   | `N/A`   |
| `username`  | user name for communication with Artifactory | `true`   | `N/A`   |

#### Copy

The following parameters are used to configure the `copy` action:

| Name        | Description                                         | Required | Default |
| ----------- | --------------------------------------------------- | -------- | ------- |
| `flat`      | enables removing source directory hierarchy         | `false`  | `false` |
| `recursive` | enables copying sub-directories for the artifact(s) | `false`  | `false` |
| `target`    | target path to copy artifact(s) to                  | `true`   | `N/A`   |


#### Delete

The following parameters are used to configure the `delete` action:

| Name        | Description                                          | Required | Default |
| ----------- | ---------------------------------------------------- | -------- | ------- |
| `recursive` | enables removing sub-directories for the artifact(s) | `false`  | `false` |

#### Docker-Promote

The following parameters are used to configure the `docker-promote` action:

| Name                     | Description                                           | Required | Default  |
| ------------------------ | ----------------------------------------------------- | -------- | -------- |
| `target_repo`            | name of the docker registry containing the image      | `true`   | `N/A`    |
| `docker_registry`        | path to image in docker registry                      | `true`   | `N/A`    |
| `target_docker_registry` | path for target image in docker registry              | `true`   | `N/A`    |
| `tag`                    | name of the tag for promoting                         | `true`   | `N/A`    |
| `target_tags`            | name of the final tags after promotion                | `true`   | `N/A`    |
| `copy`                   | set to copy instead of moving the image               | `false`  | `false`  |

#### Set-Prop

The following parameters are used to configure the `set-prop` action:

| Name    | Description                          | Required | Default |
| ------- | ------------------------------------ | -------- | ------- |
| `props` | properties to set on the artifact(s) | `true`   | `N/A`   |

#### Upload

The following parameters are used to configure the `upload` action:

| Name           | Description                                           | Required | Default |
| -------------- | ----------------------------------------------------- | -------- | ------- |
| `flat`         | enables removing source directory hierarchy           | `false`  | `false` |
| `include_dirs` | enables including sub-directories for the artifact(s) | `false`  | `false` |
| `recursive`    | enables uploading sub-directories for the artifact(s) | `false`  | `false` |
| `regexp`       | enables reading the sources as a regular expression   | `false`  | `false` |
| `sources`      | list of artifact(s) to upload                         | `true`   | `N/A`   |

## Template

COMING SOON!

## Troubleshooting

Below are a list of common problems and how to solve them:
