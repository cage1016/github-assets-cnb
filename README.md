# Github Asset Buildpack

![Version](https://img.shields.io/badge/dynamic/json?url=https://cnb-registry-api.herokuapp.com/api/v1/buildpacks/cage1016/github-assets-cnb&label=Version&query=$.latest.version)

A [Cloud Native Buildpack](https://buildpacks.io) that Download Github Assets


## Buildpack registry

https://registry.buildpacks.io/buildpacks/cage1016/github-assets-cnb

## Usage

1. Create `project.toml` if you want to embed github assets

    ```bash
    cat <<EOF >> project.toml
    # [[metadata.githubasset]]
    # repo = "GoogleContainerTools/skaffold"    # required
    # token = ""                                # optional for private repo
    # file = "skaffold-linux-amd64"             # required
    # version = ""                              # optional, default set to "latest"
    # target = "skaffold"                       # optional, default set to file name
    # untarpath = ""                            # optional, default set to asset layer (".gz", ".sz", ".xz", ".lz4", "tgz", "zip", ".tar", ".bz2")

    [[metadata.githubasset]]
    repo = "eugeneware/ffmpeg-static"
    file = "linux-x64"
    target = "ffmpeg"

    [[metadata.githubasset]]
    repo = "kkdai/youtube"
    file = "youtubedr_2.7.0_linux_arm64.tar.gz"
    untarpath = "bin"
    EOF
    ```

1. Build container image

    ```
    pack build myapp --buildpack cage1016/github-assets-cnb@2.0.0
    ```

1. Check `/layers/cage1016_github-assets-cnb`

    ![](snipaste.png)

### URI

```
urn:cnb:registry:cage1016/github-assets-cnb
```

### Supported Stacks

- google
- io.buildpacks.stacks.bionic
- io.paketo.stacks.tiny
- io.buildpacks.samples.stacks.bionic
- heroku-18
- heroku-20