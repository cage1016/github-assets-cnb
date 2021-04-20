package github

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/postal"
)

func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		entry := context.Plan.Entries[0]

		token := entry.Metadata["TOKEN"].(string)
		repo := entry.Metadata["REPO"].(string)
		version := entry.Metadata["VERSION"].(string)
		file := entry.Metadata["FILE"].(string)
		target := entry.Metadata["TARGET"].(string)

		if target == "" {
			target = file
		}

		if repo == "" || file == "" {
			fmt.Println("-----> REPO or FILE is empty, by pass")
			return packit.BuildResult{}, nil
		}

		layer, err := context.Layers.Get("github-assets")
		if err != nil {
			return packit.BuildResult{}, err
		}

		// layer.Build = true
		layer.Launch = true

		transport := cargo.NewTransport()
		service := postal.NewService(transport)
		dependency, err := service.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), "jq", "1.6", "google")
		if err != nil {
			log.Fatal(err)
		}

		// download jq for curl usage later
		fmt.Println(fmt.Sprintf("-----> Download jq %s", dependency.URI))
		dir, err := os.MkdirTemp("", "bundle")
		bundle, err := transport.Drop("", dependency.URI)
		if err != nil {
			log.Fatal("transport.Drop", err)
		}

		contents, err := io.ReadAll(bundle)
		jqPath := filepath.Join(dir, "jq")
		err = os.WriteFile(jqPath, contents, os.ModePerm)
		if err != nil {
			log.Fatal("os.WriteFile", err)
		}

		// fetch github asset id
		cmds := []string{"-H", "Accept: application/vnd.github.v3.raw"}
		if token != "" {
			cmds = append(cmds, "-H", fmt.Sprintf("Authorization: token %s", token))
		}
		cmds = append(cmds, "-s", fmt.Sprintf("https://api.github.com/repos/%s/releases", repo))

		var jqCmd string
		if version == "latest" {
			jqCmd = fmt.Sprintf(".[0].assets | map(select(.name == \"%s\"))[0].id", file)
		} else {
			jqCmd = fmt.Sprintf(". | map(select(.tag_name == \"%s\"))[0].assets | map(select(.name == \"%s\"))[0].id", version, file)
		}

		assetsId, stderr, err := Pipeline(
			exec.Command("curl", cmds...),
			exec.Command(jqPath, jqCmd),
		)
		if err != nil {
			log.Fatal(err)
		}

		// Print the stderr, if any
		if len(stderr) > 0 {
			log.Fatal(err)
		}

		downloadDir, err := ioutil.TempDir("", "downloadDir")
		if err != nil {
			return packit.BuildResult{}, err
		}
		defer os.RemoveAll(downloadDir)

		fmt.Println(fmt.Sprintf("-----> Download Github %s asset %s as %s", repo, file, target))

		// download github asset with asset id
		dcmd := []string{"-vLJo", filepath.Join(downloadDir, target), "-H", "Accept: application/octet-stream"}
		if token != "" {
			dcmd = append(dcmd, fmt.Sprintf("https://%s:@api.github.com/repos/%s/releases/assets/%s", entry.Metadata["TOKEN"], repo, string(assetsId)[:8]))
		} else {
			dcmd = append(dcmd, fmt.Sprintf("https://api.github.com/repos/%s/releases/assets/%s", repo, string(assetsId)[:8]))
		}

		err = exec.Command("curl", dcmd...).Run()
		if err != nil {
			log.Fatal("os.WriteFile", err)
		}

		launchEnvDir := filepath.Join(layer.Path, "bin")
		os.MkdirAll(launchEnvDir, os.ModePerm)

		err = exec.Command("mv", filepath.Join(downloadDir, target), launchEnvDir).Run()
		if err != nil {
			log.Fatal("os.WriteFile", err)
		}

		if err := os.Chmod(filepath.Join(launchEnvDir, target), os.ModePerm); err != nil {
			log.Fatal("os.Chmod", err)
		}

		return packit.BuildResult{
			Layers: []packit.Layer{
				layer,
			},
		}, nil
	}
}
