package main

import (
	"github.com/paketo-buildpacks/packit"

	githubassetscnb "github.com/cage1016/github-assets-cnb"
)

func main() {
	packit.Run(
		githubassetscnb.Detect(),
		githubassetscnb.Build(),
	)
}
