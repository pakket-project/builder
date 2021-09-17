package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/mholt/archiver/v3"
	"github.com/pakket-project/builder/internals/runner"
	"github.com/pakket-project/builder/util"
	"github.com/pakket-project/pakket/internals/pkg"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

type PackageInfo struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

type Info struct {
	BuildArch  string      `toml:"buildArch"`
	Repository string      `toml:"repository"`
	Package    PackageInfo `toml:"package"`
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build package_path version",
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
		err = runner.RunScript(
			path.Join(versionPath, "package"),
			"STEW_PKG_NAME="+p.Package.Name,
		)
		if err != nil {
			panic(err)
		}

		infoFile, err := os.Create(path.Join(util.TmpPkgPath, "info.toml"))
		if err != nil {
			panic(err)
		}

		info, err := toml.Marshal(Info{
			BuildArch: runtime.GOARCH,
			Package: PackageInfo{
				Name:    p.Package.Name,
				Version: version,
			},
		}) //todo: add repo
		if err != nil {
			panic(err)
		}

		_, err = infoFile.Write(info)
		if err != nil {
			panic(err)
		}

		var arch string
		if runtime.GOARCH == "arm64" {
			arch = "silicon"
		} else if runtime.GOARCH == "amd64" {
			arch = "intel"
		} else {
			panic("lol wat")
		}

		tarPath := path.Join(util.TmpPkgRootPath, fmt.Sprintf("%s-%s-%s.tar.xz", p.Package.Name, version, arch))

		err = archiver.Archive([]string{util.TmpPkgPath}, tarPath)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\n\ndone! tar: %s\npkg path: %s\n", tarPath, util.TmpPkgRootPath)
	},
}
