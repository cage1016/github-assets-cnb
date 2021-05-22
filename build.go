package cnbtest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/vacation"
)

func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		Layers := []packit.Layer{}
		var trans *Transport
		for _, entry := range context.Plan.Entries {
			tokenEnv := entry.Metadata[TokenEnv].(string)
			if tokenEnv != "" {
				if token, ok := os.LookupEnv(tokenEnv); ok {
					trans = NewTransport(SetToken(token))
				} else {
					return packit.BuildResult{}, fmt.Errorf("Repo: %s required Token ENV \"%s\" not found", entry.Metadata[Repo].(string), tokenEnv)
				}
			} else {
				trans = NewTransport()
			}

			layer, err := context.Layers.Get(strings.Replace(entry.Metadata[Repo].(string), "/", "_", 1))
			if err != nil {
				return packit.BuildResult{}, err
			}

			// layer.Build = true
			layer.Launch = true

			id, err := trans.Fetch(entry.Metadata[Repo].(string), entry.Metadata[Asset].(string), entry.Metadata[Tag].(string))
			if err != nil {
				return packit.BuildResult{}, err
			}

			bundle, mime, err := trans.Drop(entry.Metadata[Repo].(string), id)
			if err != nil {
				return packit.BuildResult{}, fmt.Errorf("failed to fetch dependency: %s", err)
			}
			defer bundle.Close()

			destination := filepath.Join(layer.Path, entry.Metadata[Destination].(string))
			stripComponents := entry.Metadata[StripComponents].(int64)

			switch mime {
			case "application/x-tar", "application/gzip", "application/x-xz":
				err = vacation.NewArchive(bundle).StripComponents(int(stripComponents)).Decompress(destination)
				if err != nil {
					return packit.BuildResult{}, err
				}
			case "application/zip":
				err = vacation.NewArchive(bundle).Decompress(destination)
				if err != nil {
					return packit.BuildResult{}, err
				}
			case "application/x-executable":
				err = Deliver(bundle, destination)
				if err != nil {
					return packit.BuildResult{}, err
				}
				err = os.Chmod(destination, os.ModePerm)
				if err != nil {
					return packit.BuildResult{}, err
				}
			default:
				err = Deliver(bundle, destination)
				if err != nil {
					return packit.BuildResult{}, err
				}
			}
			fmt.Printf("Deliver: %s -> %s\n", entry.Metadata[Asset].(string), destination)

			Layers = append(Layers, layer)
		}

		return packit.BuildResult{
			Layers: Layers,
		}, nil
	}
}

func Deliver(reader io.Reader, destination string) error {
	dir, _ := filepath.Split(destination)
	os.MkdirAll(dir, os.ModePerm)

	file, err := os.Create(destination)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}
