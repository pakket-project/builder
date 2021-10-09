package cmd

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	archiver "github.com/go-vela/archiver/v3"
	"github.com/pakket-project/builder/internals/runner"
	"github.com/pakket-project/builder/util"
	"github.com/pakket-project/pakket/internals/pkg"
	"github.com/pelletier/go-toml/v2"
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

var (
	outputDir string
)

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory")
}

var buildCmd = &cobra.Command{
	Use:   "build package_path version",
	Short: "Build packages",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if outputDir != "" {
			outputDir, err := filepath.Abs(outputDir)
			if err != nil {
				panic(err)
			}
			util.TmpRootPath = outputDir
		}

		pkgPath, err := filepath.Abs(args[0])
		if err != nil {
			panic(err)
		}

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

		versionMetadataPath := path.Join(versionPath, "metadata.toml")
		file2, err := os.ReadFile(versionMetadataPath)
		if err != nil {
			panic(err)
		}

		v, err := pkg.ParseVersion(file2)
		if err != nil {
			panic(err)
		}

		// install dependencies
		for _, dep := range v.Dependencies.Dependencies {
			var name string
			var version *string

			if strings.Contains(dep, "@") {
				splitted := strings.Split(dep, "@")
				name = splitted[0]
				version = &splitted[1]
			} else {
				name = dep
				version = nil
			}

			pkgData, err := pkg.GetPackage(name, version)
			if err != nil {
				fmt.Printf("error while installing %s: %s\n", dep, err.Error())
				continue
			}

			err = pkg.InstallPackage(*pkgData, false)
			if err != nil {
				fmt.Printf("error while installing %s: %s\n", dep, err.Error())
				continue
			}
		}

		// install build dependencies
		for _, dep := range v.Dependencies.BuildDependencies {
			var name string
			var version *string

			if strings.Contains(dep, "@") {
				splitted := strings.Split(dep, "@")
				name = splitted[0]
				version = &splitted[1]
			} else {
				name = dep
				version = nil
			}

			pkgData, err := pkg.GetPackage(name, version)
			if err != nil {
				fmt.Printf("error while installing %s: %s\n", dep, err.Error())
				continue
			}

			err = pkg.InstallPackage(*pkgData, false)
			if err != nil {
				fmt.Printf("error while installing %s: %s\n", dep, err.Error())
				continue
			}
		}

		err = util.CreateTempFolder(p.Package.Name)
		if err != nil {
			panic(err)
		}

		err = util.DownloadArchive(v.Url, util.TmpSrcPath)
		if err != nil {
			panic(err)
		}

		os.Chdir(util.TmpSrcPath)
		err = runner.RunScript(
			path.Join(versionPath, "package"),
			"PAKKET_PKG_NAME="+p.Package.Name,
		)
		if err != nil {
			fmt.Printf("error running script: %s\n", err.Error())
			os.Exit(1)
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

		tarPath := path.Join(util.TmpPath, fmt.Sprintf("%s-%s-%s.tar.xz", p.Package.Name, version, arch))

		err = archiver.Archive([]string{util.TmpPkgPath}, tarPath)
		if err != nil {
			fmt.Printf("error creating archive %s: %s\n", tarPath, err.Error())
			os.Exit(1)
		}

		err = os.RemoveAll(util.TmpSrcPath)
		if err != nil {
			panic(err)
		}

		user, err := user.Lookup(os.Getenv("SUDO_USER"))
		if err != nil {
			panic(err)
		}
		uid, err := strconv.Atoi(user.Uid)
		if err != nil {
			panic(err)
		}
		gid, err := strconv.Atoi(user.Gid)
		if err != nil {
			panic(err)
		}

		err = filepath.WalkDir(util.TmpRootPath, func(path string, d fs.DirEntry, err error) error {
			return os.Chown(path, uid, gid)
		})
		if err != nil {
			panic(err)
		}

		// update checksums
		tarData, err := os.ReadFile(tarPath)
		if err != nil {
			panic(err)
		}
		checksum := fmt.Sprintf("%x", sha256.Sum256(tarData))

		versionContent := strings.Split(string(file2), "\n")

		// TODO: fix this
		done := false
		for i, line := range versionContent {
			if strings.HasPrefix(line, "checksum") {
				if strings.Contains(versionContent[i-1], arch) {
					versionContent[i] = "checksum = '" + checksum + "'"
					done = true
				}
			}
		}

		err = os.WriteFile(versionMetadataPath, []byte(strings.Join(versionContent, "\n")), 0666)
		if err != nil {
			panic(err)
		}

		if !done {
			fmt.Printf("tar: %s\n", tarPath)
			fmt.Printf("pkg path: %s\n", util.TmpPkgPath)
			fmt.Printf("checksum: %s\n", checksum)
			panic("ERROR: failed to update metadata.toml with new hash!")
		}

		fmt.Println("done!")
		fmt.Printf("tar: %s\n", tarPath)
		fmt.Printf("pkg path: %s\n", util.TmpPkgPath)
		fmt.Printf("checksum: %s\n", checksum)
	},
}
