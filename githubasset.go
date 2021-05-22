package cnbtest

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Githubasset is a representation of a buildpack githubasset.
type Githubasset struct {
	Repo string `toml:"repo"`

	Asset string `toml:"asset"`

	Tag string `toml:"tag"`

	TokenEnv string `toml:"token_env"`

	Destination string `toml:"destination"`

	StripComponents int64 `toml:"strip_components"`
}

func parseBuildpack(path string) ([]Githubasset, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse buildpack.toml: %w", err)
	}

	var buildpack struct {
		Metadata struct {
			Dependencies []Githubasset `toml:"githubassets"`
		} `toml:"metadata"`
	}
	_, err = toml.DecodeReader(file, &buildpack)
	if err != nil {
		return nil, fmt.Errorf("failed to parse buildpack.toml: %w", err)
	}

	return buildpack.Metadata.Dependencies, nil
}

func stacksInclude(stacks []string, stack string) bool {
	for _, s := range stacks {
		if s == stack {
			return true
		}
	}
	return false
}
