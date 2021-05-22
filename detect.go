package cnbtest

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		plan := packit.DetectResult{
			Plan: packit.BuildPlan{},
		}

		var githubassets []Githubasset
		var err error
		if p, ok := os.LookupEnv("GITHUB_ASSETS_PATH"); ok {
			githubassets, err = parseBuildpack(filepath.Join(context.WorkingDir, p))
			if err != nil {
				return packit.DetectResult{}, err
			}
		} else {
			githubassets, err = parseBuildpack(filepath.Join(context.WorkingDir, "project.toml"))
			if err != nil {
				return packit.DetectResult{}, nil
			}
		}

		if len(githubassets) == 0 {
			return plan, nil
		}

		plan.Plan.Provides = []packit.BuildPlanProvision{}
		plan.Plan.Requires = []packit.BuildPlanRequirement{}

		for _, asset := range githubassets {
			name := strings.Replace(asset.Repo, "/", "_", 1)
			plan.Plan.Provides = append(plan.Plan.Provides, packit.BuildPlanProvision{
				Name: name,
			})
			plan.Plan.Requires = append(plan.Plan.Requires, packit.BuildPlanRequirement{
				Name: name,
				Metadata: map[string]interface{}{
					Repo:            asset.Repo,
					Asset:           asset.Asset,
					Tag:             asset.Tag,
					TokenEnv:        asset.TokenEnv,
					Destination:     asset.Destination,
					StripComponents: asset.StripComponents,
				},
			})
		}
		return plan, nil
	}
}
