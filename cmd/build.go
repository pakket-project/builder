package cmd

import (
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/mholt/archiver/v3"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/stewproject/builder/internals/runner"
	"github.com/stewproject/builder/util"
	"github.com/stewproject/stew/internals/pkg"
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
	Use:   "build package version GH_org",
	Short: "Build packages",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		pkgPath := args[0]
		version := args[1]
		GH_org := args[2]
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

		tarPath := path.Join(util.TmpPkgRootPath, p.Package.Name+".tar.xz")

		err = archiver.Archive([]string{util.TmpPkgPath}, tarPath)
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

		uploadCmd := exec.Command("skopeo", "copy", "tarball:"+tarPath, "docker://ghcr.io/"+GH_org+"/packages/"+p.Package.Name+":"+version+"_"+arch)

		uploadCmd.Stderr = os.Stderr
		uploadCmd.Stdin = os.Stdin
		uploadCmd.Stdout = os.Stdout

		err = uploadCmd.Run()
		if err != nil {
			panic(err)
		}
	},
}
