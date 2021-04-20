package github

import (
	"os"

	"github.com/paketo-buildpacks/packit"
)

var md = map[string]string{
	"TOKEN":   "",
	"REPO":    "",
	"VERSION": "latest",
	"FILE":    "",
	"TARGET":  "",
}

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		for k, _ := range md {
			if e, ok := os.LookupEnv(k); ok {
				md[k] = e
			}
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "github-assets"},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name:     "github-assets",
						Metadata: md,
					},
				},
			},
		}, nil
	}
}
