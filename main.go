package estailwind

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

const tailwindCmd = "tailwindcss"
const tailwindModulesPath = "node_modules/.bin/" + tailwindCmd
const tailwindConfig = "tailwind.config.js"

func findFile(basePath string, fileName string) (string, error) {
	base := filepath.Dir(basePath)
	filePath := base

	found := false
	// no point in going to the root, it won't be there
	for base != "/" {
		filePath = filepath.Join(base, fileName)

		// found it!!!
		if _, err := os.Stat(filePath); err == nil {
			found = true
			break
		}

		// go a parent level, maybe it is there
		base = filepath.Join(base, "..")
	}

	if found {
		return filePath, nil
	}

	return "", nil
}

func findTailwindCmd(baseCmdPath string) (string, error) {
	tailwindPath, err := findFile(baseCmdPath, tailwindModulesPath)
	if err != nil {
		return "", err
	}

	if tailwindPath != "" {
		return tailwindPath, nil
	}

	tailwindPath, err = findFile(baseCmdPath, tailwindCmd)
	if err != nil {
		return "", err
	}

	if tailwindPath != "" {
		return tailwindPath, nil
	}

	return exec.LookPath(tailwindCmd)
}

func runTailwind(inputFile string, outputFile string, isMinify bool) error {
	cmdPath, err := findTailwindCmd(inputFile)
	if err != nil {
		return err
	}

	tailwindConfig, err := findFile(inputFile, tailwindConfig)
	if err != nil {
		return err
	}

	args := []string{
		"-i", inputFile,
		"-o", outputFile,
	}
	if tailwindConfig != "" {
		args = append(args, "-c", tailwindConfig)
	}

	if isMinify {
		args = append(args, "-m")
	}

	cmd := exec.Command(cmdPath, args...)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func NewTailwindPlugin(shouldMinify bool) api.Plugin {
	return api.Plugin{
		Name: "tailwind",
		Setup: func(b api.PluginBuild) {
			tmpFiles := []string{}

			b.OnResolve(api.OnResolveOptions{
				Filter:    `.\.(css)$`,
				Namespace: "file",
			}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				sourceFullPath := filepath.Join(args.ResolveDir, args.Path)
				source, err := os.ReadFile(sourceFullPath)
				if err != nil {
					return api.OnResolveResult{Path: sourceFullPath}, err
				}

				if !strings.Contains(string(source), "@tailwind") {
					return api.OnResolveResult{Path: sourceFullPath}, nil
				}

				tmpFile := strings.ReplaceAll(
					filepath.Base(sourceFullPath),
					filepath.Ext(sourceFullPath),
					"") + ".tmp.css"
				tmpFilePath := filepath.Join(filepath.Dir(sourceFullPath), tmpFile)
				tmpFiles = append(tmpFiles, tmpFilePath)

				err = runTailwind(sourceFullPath, tmpFilePath, shouldMinify)
				return api.OnResolveResult{
					Path: tmpFilePath,
				}, err
			})
			b.OnEnd(func(result *api.BuildResult) (api.OnEndResult, error) {
				// remove tmp files
				for _, tmp := range tmpFiles {
					_ = os.Remove(tmp)
				}

				return api.OnEndResult{
					Errors:   result.Errors,
					Warnings: result.Warnings,
				}, nil
			})
		},
	}
}
