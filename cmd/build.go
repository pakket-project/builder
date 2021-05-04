package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/stewproject/builder/internals/runner"
	"github.com/stewproject/builder/util"
	"github.com/stewproject/stew/internals/pkg"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build package version",
	Short: "Build packages",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pkgPath := args[0]
		version := args[1]
		versionPath := path.Join(pkgPath, version)

		file, err := os.ReadFile(path.Join(pkgPath, "package.toml"))
		if err != nil {
			panic(err)
		}

		p, err := pkg.ParsePackage(file)
		if err != nil {
			panic(err)
		}

		file2, err := os.ReadFile(path.Join(versionPath, "metadata.toml"))
		if err != nil {
			panic(err)
		}

		v, err := pkg.ParseVersion(file2)
		if err != nil {
			panic(err)
		}

		err = util.CreateTempFolder(p.Package.Name)
		if err != nil {
			panic(err)
		}

		err = util.DownloadArchive(v.Url, util.TmpSrcPath)
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(util.TmpSrcPath)

		os.Chdir(util.TmpSrcPath)
		if err = runner.RunScript(
			path.Join(versionPath, "package"),
			"STEW_PKG_PATH="+util.TmpPkgPath,
			"STEW_SRC_DIR="+util.TmpSrcPath,
			"STEW_PKG_NAME="+p.Package.Name,
		); err != nil {
			panic(err)
		}
	},
}
