package github

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
)

type Item struct {
	Repo      string `json:"repo"`
	File      string `json:"file"`
	Target    string `json:"target"`
	Version   string `json:"version"`
	Token     string `json:"token"`
	UnTarPath string `json:"untarpath"`
}

func (p *Item) UnmarshalJSON(data []byte) error {
	type Alias Item

	pr := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &pr); err != nil {
		return nil
	}

	if p.Version == "" {
		p.Version = "latest"
	}

	return nil
}

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		parser := cargo.NewBuildpackParser()

		config, err := parser.Parse(filepath.Join(context.WorkingDir, "project.toml"))
		if err != nil {
			return packit.DetectResult{}, nil
		}

		var items []Item
		if ga, ok := config.Metadata.Unstructured["githubasset"]; ok {
			err := json.Unmarshal(ga.(json.RawMessage), &items)
			if err != nil {
				return packit.DetectResult{}, nil
			}
		}

		plan := packit.DetectResult{Plan: packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{},
		}}

		for _, v := range items {
			name := strings.Replace(v.Repo, "/", "_", 1)
			plan.Plan.Provides = append(plan.Plan.Provides, packit.BuildPlanProvision{
				Name: name,
			})
			plan.Plan.Requires = append(plan.Plan.Requires, packit.BuildPlanRequirement{
				Name: name,
				Metadata: map[string]string{
					REPO:      v.Repo,
					FILE:      v.File,
					TARGET:    v.Target,
					VERSION:   v.Version,
					TOKEN:     v.Token,
					UNTARPATH: v.UnTarPath,
				},
			})
		}

		return plan, nil
	}
}
