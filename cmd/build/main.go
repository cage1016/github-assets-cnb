package main

import (
	"github.com/paketo-buildpacks/packit"

	"github.com/cage1016/github-assets-cnb/github"
)

func main() {
	packit.Build(github.Build())
}
