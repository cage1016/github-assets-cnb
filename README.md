# Github Asset Buildpack

![Version](https://img.shields.io/badge/dynamic/json?url=https://cnb-registry-api.herokuapp.com/api/v1/buildpacks/cage1016/github-assets-cnb&label=Version&query=$.latest.version)

A [Cloud Native Buildpack](https://buildpacks.io) that Download Github Assets


## Buildpack registry

https://registry.buildpacks.io/buildpacks/cage1016/github-assets-cnb

## Usage

1. Create `project.toml` if you want to embed github assets

    ```bash
    cat <<EOF >> project.toml
    # [[build.env]]
    # optional, github token for private assets 
    # name = "TOKEN"
    # value = "<github-token>"

    # skaffold
    [[build.env]]
    # required
    name = "REPO"
    value = "GoogleContainerTools/skaffold"

    [[build.env]]
    # required
    name = "FILE"
    value = "skaffold-linux-amd64"

    [[build.env]]
    # optional, default set to FILE value
    name = "TARGET"
    value = "skaffold"

    [[build.env]]
    # optional, default set to 'latest'
    name = "VERSION"
    value = "v1.22.0"
    EOF
    ```


```
pack build myapp --buildpack cage1016/github-assets-cnb@1.0.0
```

### URI

```
urn:cnb:registry:cage1016/github-assets-cnb
```

### Supported Stacks

- google
- io.buildpacks.stacks.bionic
- io.buildpacks.samples.stacks.bionic
- heroku-18
- heroku-20